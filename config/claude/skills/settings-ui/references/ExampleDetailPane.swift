//
//  ExampleDetailPane.swift
//  {{AppName}}
//
//  Template showing common control patterns inside a settings detail pane.
//  Every pane follows the same structure: Form > Section > controls,
//  with .formStyle(.grouped) + .scrollContentBackground(.hidden).
//

import SwiftUI

// MARK: - General Settings Pane (example with toggles, buttons, labeled content)

struct GeneralSettingsPane: View {
    @AppStorage("launchAtLogin") private var launchAtLogin = false
    @AppStorage("exportDirectoryPath") private var exportDirectoryPath = ""

    var body: some View {
        Form {
            Section("Save Location") {
                LabeledContent("Export folder") {
                    HStack(spacing: 8) {
                        Image(systemName: "folder.fill")
                            .foregroundStyle(.blue)
                            .font(.system(size: 14))

                        Text(exportDirectoryPath.isEmpty ? "~/Documents" : exportDirectoryPath)
                            .font(.system(size: 13))
                            .lineLimit(1)
                            .truncationMode(.middle)
                            .foregroundStyle(.primary)
                    }
                }

                HStack(spacing: 8) {
                    Button("Choose Folder...") {
                        // Open NSOpenPanel here
                    }
                    .controlSize(.small)

                    Button("Use Default") {
                        exportDirectoryPath = ""
                    }
                    .controlSize(.small)
                    .disabled(exportDirectoryPath.isEmpty)
                }
            }

            Section("System") {
                // Toggle with description text
                Toggle(isOn: $launchAtLogin) {
                    VStack(alignment: .leading, spacing: 2) {
                        Text("Launch at Login")
                        Text("Start the app automatically when you sign in.")
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
                .toggleStyle(.switch)
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
        .contentMargins(.top, 8, for: .scrollContent)
    }
}

// MARK: - Appearance Settings Pane (example with segmented picker)

struct AppearanceSettingsPane: View {
    @AppStorage("appTheme") private var appTheme = "system"

    var body: some View {
        Form {
            Section("Theme") {
                Picker("Appearance", selection: $appTheme) {
                    Text("System").tag("system")
                    Text("Light").tag("light")
                    Text("Dark").tag("dark")
                }
                .pickerStyle(.segmented)
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
        .contentMargins(.top, 8, for: .scrollContent)
    }
}

// MARK: - About Settings Pane (example with app identity + links)

struct AboutSettingsPane: View {
    private var versionText: String {
        let version = Bundle.main.infoDictionary?["CFBundleShortVersionString"] as? String
        let build = Bundle.main.infoDictionary?["CFBundleVersion"] as? String
        switch (version, build) {
        case let (v?, b?): return "Version \(v) (\(b))"
        case let (v?, nil): return "Version \(v)"
        default: return "Version 1.0"
        }
    }

    var body: some View {
        Form {
            Section {
                HStack(alignment: .center, spacing: 16) {
                    Image(nsImage: NSApp.applicationIconImage)
                        .resizable()
                        .aspectRatio(contentMode: .fit)
                        .frame(width: 72, height: 72)

                    VStack(alignment: .leading, spacing: 6) {
                        Text("MyApp")
                            .font(.largeTitle.bold())

                        Text(versionText)
                            .font(.subheadline)
                            .foregroundStyle(.secondary)

                        Text("A brief description of your app.")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                    }
                }
            }

            Section("Project") {
                Text("A longer description of what the app does and why it exists.")
                    .foregroundStyle(.secondary)

                Link("GitHub", destination: URL(string: "https://github.com/you/your-app")!)
            }

            Section("Credits") {
                Text("Built by Your Name")
                    .foregroundStyle(.secondary)
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
        .contentMargins(.top, 8, for: .scrollContent)
    }
}

// MARK: - Advanced Example: Slider + Color Picker + Dropdown Picker

struct AdvancedSettingsPane: View {
    @AppStorage("indicatorSize") private var indicatorSize = 48.0
    @AppStorage("enableFeature") private var enableFeature = true
    @AppStorage("outputFormat") private var outputFormat = "png"

    var body: some View {
        Form {
            Section("Features") {
                Toggle(isOn: $enableFeature) {
                    VStack(alignment: .leading, spacing: 2) {
                        Text("Enable experimental feature")
                        Text("This feature is still in development.")
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
                .toggleStyle(.switch)
            }

            Section("Output") {
                // Dropdown picker
                Picker("Format", selection: $outputFormat) {
                    Text("PNG").tag("png")
                    Text("JPEG").tag("jpeg")
                    Text("HEIC").tag("heic")
                }
                .pickerStyle(.menu)
            }

            Section("Indicator") {
                // Slider with value readout
                LabeledContent("Size") {
                    HStack(spacing: 12) {
                        Slider(value: $indicatorSize, in: 24...96, step: 2)
                            .frame(width: 180)

                        Text("\(Int(indicatorSize)) pt")
                            .monospacedDigit()
                            .foregroundStyle(.secondary)
                            .frame(width: 46, alignment: .trailing)
                    }
                }
            }
        }
        .formStyle(.grouped)
        .scrollContentBackground(.hidden)
        .contentMargins(.top, 8, for: .scrollContent)
    }
}
