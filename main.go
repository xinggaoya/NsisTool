package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"text/template"
)

const nsisTemplate = `!include "MUI2.nsh"

OutFile "{{.Name}}-installer.exe"

Name "{{.Name}}"
Caption "{{.Name}} Installer"
BrandingText "{{.Company}}"

InstallDir "$PROGRAMFILES\{{.Name}}"
InstallDirRegKey HKCU "Software\{{.Name}}" "Install_Dir"

!define MUI_ABORTWARNING
!define MUI_ICON "{{.IconPath}}"
!define MUI_UNICON "{{.IconPath}}"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

VIProductVersion "{{.Version}}"
VIAddVersionKey /LANG=1033 "ProductName" "{{.Name}}"
VIAddVersionKey /LANG=1033 "CompanyName" "{{.Company}}"
VIAddVersionKey /LANG=1033 "FileDescription" "{{.Name}} Installer"
VIAddVersionKey /LANG=1033 "FileVersion" "{{.Version}}"
VIAddVersionKey /LANG=1033 "ProductVersion" "{{.Version}}"
VIAddVersionKey /LANG=1033 "OriginalFilename" "{{.Name}}-installer.exe"
VIAddVersionKey /LANG=1033 "InternalName" "{{.Name}}Installer"
VIAddVersionKey /LANG=1033 "LegalCopyright" "© {{.Company}}. All rights reserved."

Section "Install"

    CreateDirectory "$INSTDIR"
    SetOutPath "$INSTDIR"
    File "{{.ProgramPath}}"
    CreateShortcut "$DESKTOP\{{.Name}}.lnk" "$INSTDIR\{{.Name}}.exe"
    WriteRegStr HKCU "Software\{{.Name}}" "Install_Dir" "$INSTDIR"
    WriteUninstaller "$INSTDIR\Uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "DisplayName" "{{.Name}}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "UninstallString" "$INSTDIR\Uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "InstallLocation" "$INSTDIR"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "DisplayIcon" "$INSTDIR\{{.Name}}.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "DisplayVersion" "{{.Version}}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "Publisher" "{{.Company}}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "VersionMajor" "1"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "VersionMinor" "0"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "Version" "{{.Version}}"
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "NoModify" 1
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}" "NoRepair" 1

SectionEnd

Section "Uninstall"

    Delete "$INSTDIR\{{.Name}}.exe"
    Delete "$DESKTOP\{{.Name}}.lnk"
    RMDir /r "$INSTDIR"
    DeleteRegKey HKCU "Software\{{.Name}}"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\{{.Name}}"

SectionEnd
`

type Program struct {
	Name        string `json:"name"`
	Company     string `json:"company"`
	Version     string `json:"version"`
	IconPath    string `json:"icon_path"`
	ProgramPath string `json:"program_path"`
}

func generateNsisScript(program Program) (string, error) {
	_, err := template.New("nsis").Parse(nsisTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing NSIS template: %w", err)
	}

	builder := &template.Template{}
	builder = template.Must(builder.Parse(nsisTemplate))

	var scriptBuffer []byte
	buffer := bytes.NewBuffer(scriptBuffer)
	err = builder.Execute(buffer, program)
	if err != nil {
		return "", fmt.Errorf("error executing NSIS template: %w", err)
	}

	return buffer.String(), nil
}

func run() {
	jsonFile, err := os.Open("nsis.json")
	if err != nil {
		slog.Error("Error opening nsis.json", "error", err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var program Program
	json.Unmarshal(byteValue, &program)

	script, _ := generateNsisScript(program)
	scriptFilename := fmt.Sprintf("%s_installer.nsi", program.Name)
	err = os.WriteFile(scriptFilename, []byte(script), 0644)
	if err != nil {
		slog.Error("Error writing NSIS script", "error", err)
		return
	}

	// 删除
	defer os.Remove(scriptFilename)

	cmd := exec.Command("makensis", scriptFilename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		slog.Error("Error running makensis", "error", err)
		return
	}
	slog.Info("Successfully created installer", "installer", program.ProgramPath)
}

func main() {
	run()
}
