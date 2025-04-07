// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	backup "github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func exportSystemdConfiguration(outputDir string, dryRun bool) error {
	filesToBackup := gatherSystemdItems()

	if dryRun {
		log.Info().Msgf(L("Would backup %s"), filesToBackup)
		return nil
	}
	// Create output file
	out, err := os.Create(path.Join(outputDir, backup.SystemdConfBackupFile))
	if err != nil {
		return fmt.Errorf(L("failed to create Systemd backup tarball: %w"), err)
	}
	defer out.Close()

	// Prepare tar buffer
	tw := tar.NewWriter(out)
	defer tw.Close()

	for _, fileToBackup := range filesToBackup {
		f, err := os.Open(fileToBackup)
		if err != nil {
			return err
		}
		fstat, _ := f.Stat()
		h, err := tar.FileInfoHeader(fstat, "")
		if err != nil {
			return err
		}
		// Produced header does not have full path, overwrite it
		h.Name = fileToBackup
		if fstat.IsDir() {
			h.Name += "/"
		}
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		if fstat.IsDir() {
			continue
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
	}
	return nil
}

// For each container get service file, service.d and its content.
func gatherSystemdItems() []string {
	result := []string{}

	for _, service := range utils.UyuniServices {
		serviceName, skip := findService(service.Name)
		if skip {
			continue
		}

		result = append(result, podman.GetServicePath(serviceName))
		// For single mandatory replica following returns 0 so loop is skipped
		for i := 0; i < systemd.CurrentReplicaCount(service.Name); i++ {
			result = append(result, podman.GetServicePath(fmt.Sprintf("%s%d", serviceName, i)))
		}
		serviceConfDir := podman.GetServiceConfFolder(serviceName)
		serviceFiles, err := os.ReadDir(serviceConfDir)
		if err != nil {
			log.Debug().Msgf("Service configuration directory %s not found, skipping", serviceConfDir)
			continue
		}
		result = append(result, serviceConfDir)
		for _, entry := range serviceFiles {
			result = append(result, path.Join(serviceConfDir, entry.Name()))
		}
	}
	return result
}

func findService(name string) (serviceName string, skip bool) {
	skip = false
	serviceName = name
	if !systemd.HasService(serviceName) {
		// with optional or more replicas we have service template, check if the service exists at all
		serviceName = name + "@"
		if !systemd.HasService(serviceName) {
			log.Debug().Msgf("No service found for %s, skipping", name)
			skip = true
		}
	}
	return
}
