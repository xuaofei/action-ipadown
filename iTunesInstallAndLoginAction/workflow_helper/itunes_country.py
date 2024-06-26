import subprocess
import time
from pywinauto.application import Application
from win32con import *
import sys
import requests
import json


print("Launching iTunes...")
webAddress = "http://192.168.0.210"

def initITunes():
    subprocess.call('taskkill /f /im APSDaemon*', shell=True)
    subprocess.call('taskkill /f /im iTunes*', shell=True)

    app = Application().start(r"C:\Program Files\iTunes\iTunes.exe")
    app.wait_cpu_usage_lower(50)
    time.sleep(8)

    def debugTopWin():
        topwin = app.top_window().wait('exists')
        texts = []
        texts += topwin.texts()
        print("-- Cur top win: %s, texts: %s" % (topwin, topwin.texts()))
        print("-- Cur top win: %s, window_text: %s" % (topwin, topwin.window_text()))
        for c in topwin.iter_children():
            texts += c.texts()
            print("-- Cur top win: %s, texts: %s" % (topwin, c.texts()))
            print("-- Cur top win: %s, window_text: %s" % (topwin, c.window_text()))
        print("-----------------------------------------------------------------")
        # print("-- Cur top win: %s, texts: %s" % (topwin, texts))
        return "-- Cur top win: %s, texts: %s" % (topwin, texts)

    def cleanAllDialog():
        while True:
            topwin = app.top_window().wait('exists')
            if 'Dialog' in topwin.class_name():
                print("    Closing dialog %s" % topwin.window_text())
                app.top_window().Button0.click()
            elif 'Tour' in topwin.window_text():
                print("    Closing Window %s" % topwin.window_text())
                topwin.close()
            else:
                break
            
            app.wait_cpu_usage_lower(50)
            time.sleep(5)

    def cleanFailedDialog():
        while True:
            topwin = app.top_window().wait('exists')
            # dialogWrap = topwin.wait('ready')
            # print("    friendly_class_name %s" % dialogWrap.friendly_class_name())
            print("    class_name %s" % topwin.class_name())
            print("    friendly_class_name %s" % topwin.friendly_class_name())
            print("    window_text %s" % topwin.window_text())

            if 'Dialog' in topwin.friendly_class_name() and 'Verification Failed' in topwin.window_text():
                print("    Closing dialog %s" % topwin.window_text())
                app.top_window().Button0.click()
            else:
                break

            app.wait_cpu_usage_lower(50)
            time.sleep(5)

    # Click all first-time dialogs (like License Agreements, missing audios)
    # cleanAllDialog()

    # Calm down a bit before main window operations
    app.wait_cpu_usage_lower(50)

    for i in range(10):
        debugTopWin()
        time.sleep(3)

    print("-------cleanFailedDialog------")
    cleanFailedDialog()


for init_i in range(3):
    try:
        initITunes()
        break
    except Exception as e:
        print("Init iTunes %d: Failed with %s" % (init_i, e))
        import traceback; traceback.print_exc()
        time.sleep(8)

print("Init iTunes Successfully!")