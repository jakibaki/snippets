# just a quick demo of the adb fastboot thing, untested

from adb import fastboot

dev = fastboot.FastbootCommands()
dev.ConnectDevice()

# dev.Oem("unlock")
dev.RebootBootloader()

dev.Close()