import Cocoa

struct FocusedWindow : Equatable, Codable {
    let appName: String
    let windowTitle: String
    let bundleIdentifier: String
}

class EventsLoggerAppDelegate : NSObject, NSApplicationDelegate {

    var lastFocusedWindow: FocusedWindow = FocusedWindow(appName: "", windowTitle: "", bundleIdentifier: "")
    
    func applicationDidFinishLaunching(_ notification: Notification) {
        handle()
        NSWorkspace.shared.notificationCenter.addObserver(self,
            selector: #selector(activatedApp),
            name: NSWorkspace.didActivateApplicationNotification,
            object: nil)
    }

    @objc dynamic func activatedApp(_ notification: Notification) {
        handle()
    }

    func handle() -> Void {
        let pid = NSWorkspace.shared.frontmostApplication!.processIdentifier
        let appRef = AXUIElementCreateApplication(pid)
        var value: AnyObject?
        AXUIElementCopyAttributeValue(appRef, kAXWindowsAttribute as CFString, &value)
        // the first window of the front most application is the front most window
        if let targetWindow = (value as? [AXUIElement])?.first {
            var title: AnyObject?
            AXUIElementCopyAttributeValue(targetWindow, kAXTitleAttribute as CFString, &title)
            if let windowTitle = title as? String {
                let currentWindow = FocusedWindow(appName: NSWorkspace.shared.frontmostApplication?.localizedName ?? "", windowTitle: windowTitle, bundleIdentifier: NSWorkspace.shared.frontmostApplication?.bundleIdentifier ?? "")
                if currentWindow != lastFocusedWindow {
                    lastFocusedWindow = currentWindow
                    storeEventToServer(window: lastFocusedWindow)
                }
            }
        }
    }

    func storeEventToServer(window: FocusedWindow) -> Void {
        guard let url = URL(string: "http://localhost:6969/event") else {
            print("Invalid URL")
            return
        }
        do {
            let encoder = JSONEncoder()
            encoder.keyEncodingStrategy = .convertToSnakeCase
            let jsonData = try encoder.encode(window)
            // debug
            // print(String(data: jsonData, encoding: .utf8)!)
            var request = URLRequest(url: url)
            request.httpMethod = "POST"
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = jsonData
            let task = URLSession.shared.dataTask(with: request) { _, response, error in
                if let error = error {
                    print("Error making the request: \(error.localizedDescription)")
                    return
                }

                // Handle the response
                if let response = response as? HTTPURLResponse {
                    if response.statusCode != 200 {
                        print("Error: HTTP status code \(response.statusCode)")
                    }
                }
            }

            task.resume()
        } catch {
            print("Failed to make request: \(error.localizedDescription)")
        }
    }

    func applicationWillTerminate(_ notification: Notification) {
        // Clean up observers
        NSWorkspace.shared.notificationCenter.removeObserver(self)
    }
}

func main() {
    setbuf(stdout, nil)
    let delegate = EventsLoggerAppDelegate()
    NSApplication.shared.delegate = delegate
    NSApplication.shared.run()
}

main()
