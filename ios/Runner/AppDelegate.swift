import UIKit
import Flutter
import Core // 引入 Go 编译的库

@UIApplicationMain
@objc class AppDelegate: FlutterAppDelegate {
    override func application(
        _ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
    ) -> Bool {
        let controller = window?.rootViewController as! FlutterViewController
        let channel = FlutterMethodChannel(name: "com.torrent.app/core", binaryMessenger: controller.binaryMessenger)
        
        channel.setMethodCallHandler({ (call, result) in
            switch call.method {
            case "InitEngine":
                let dir = call.arguments as! String
                CoreInitEngine(dir)
                result(nil)
            case "AddMagnet":
                let magnet = call.arguments as! String
                result(CoreAddMagnet(magnet))
            case "GetFileList":
                result(CoreGetFileList())
            case "StartSequentialDownload":
                let index = call.arguments as! Int
                result(CoreStartSequentialDownload(Int32(index)))
            case "GetDownloadStats":
                result(CoreGetDownloadStats())
            case "GetFilePath":
                let index = call.arguments as! Int
                result(CoreGetFilePath(Int32(index)))
            default:
                result(FlutterMethodNotImplemented)
            }
        })
        
        GeneratedPluginRegistrant.register(with: self)
        return super.application(application, didFinishLaunchingWithOptions: launchOptions)
    }
}