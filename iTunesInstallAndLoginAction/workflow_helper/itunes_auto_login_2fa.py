import subprocess
import time
from pywinauto.application import Application
from win32con import *
import sys
import requests
import json
import logging
from rich.logging import RichHandler
from rich.console import Console
import rich
import error_code

rich.get_console().file = sys.stderr
if rich.get_console().width < 100:
    rich.get_console().width = 100

logging_handler = RichHandler(rich_tracebacks=True)
logging.basicConfig(
    level="INFO",
    format="%(message)s",
    datefmt="[%X]",
    handlers=[logging_handler]
)
logging.getLogger('urllib3').setLevel(logging.WARNING)
retryLogger = logging.getLogger('urllib3.util.retry')
retryLogger.setLevel(logging.DEBUG)
retryLogger.handlers = [logging_handler]
retryLogger.propagate = False

logger = logging.getLogger('main')


logger.info("Launching iTunes...")
webAddress = sys.argv[1]
taskId = sys.argv[2]

cpu_usage = 30
app = None

def reportResult(code, msg):
    if msg == None:
        msg = ""

    logger.info('reportResult code:%d msg:%s' % (code, msg))
    # 上报到服务端
    data_json = json.dumps({'task_id': taskId, "code": code, "msg": msg})
    url = webAddress + '/scriptReportResult'
    responseData = requests.post(url, data_json)
    logger.info("reportResult status_code:%d " % responseData.status_code)


def debugTopWin():
    topwin = app.top_window().wait('exists')
    texts = []

    texts += topwin.texts()
    for c in topwin.iter_children():
        texts += c.texts()
    logger.info("-- Cur top win: %s, texts: %s" % (topwin, texts))
    return "-- Cur top win: %s, texts: %s" % (topwin, texts)

def cleanAllDialog():
    while True:
        topwin = app.top_window().wait('exists')
        if 'Dialog' in topwin.class_name():
            logger.info("    Closing dialog %s" % topwin.window_text())
            app.top_window().Button0.click()
        elif 'Tour' in topwin.window_text():
            logger.info("    Closing Window %s" % topwin.window_text())
            topwin.close()
        else:
            break

        app.wait_cpu_usage_lower(cpu_usage)
        time.sleep(5)

def cleanFailedDialog():
    while True:
        topwin = app.top_window().wait('exists')
        if 'Dialog' in topwin.friendly_class_name() and 'Verification Failed' in topwin.window_text():
            print("    Closing dialog %s" % topwin.window_text())
            app.top_window().Button0.click()
        else:
            break

        app.wait_cpu_usage_lower(cpu_usage)
        time.sleep(5)

def cleanWelcomeDialog():
    # Click main window's first-time question ("No thanks" button)
    try:
        buttonText = app.iTunes.Button11.wait('ready').window_text()
        logger.info('Button11 text is: %s' % buttonText)
        if 'Search' not in buttonText:
            logger.info("Clicked 'No Thanks' Button!")
            app.iTunes.Button11.click_input()
            app.wait_cpu_usage_lower(cpu_usage)
            time.sleep(4)
        else:
            raise Exception('stub')
    except:
        logger.info("Not founding 'No Thanks' Button, passing on...")


def loginItunes():
    # Start logging in by clicking toolbar menu "Account"
    logger.info("Clicking Account menu...")
    app.iTunes.Application.Static3.click()
    app.wait_cpu_usage_lower(cpu_usage)
    time.sleep(3)

    debugTopWin()

    # Detect whether we have "&S" in popup, which refers to "Sign in"
    popup = app.PopupMenu
    if '&S' not in popup.menu().item(1).text():
        popup.close()
        logger.info("Already logged in!")
        return

    logger.info("Signin menu presented, clicking to login!")
    # not log in
    popup.menu().item(1).click_input()
    app.wait_cpu_usage_lower(cpu_usage)
    time.sleep(8)
    debugTopWin()

    for i in range(15):
        dialog = app.top_window()
        dialogWrap = dialog.wait('ready')
        assert dialogWrap.friendly_class_name() == 'Dialog'
        logger.info("friendly_class_name is %s" % dialogWrap.friendly_class_name())
        time.sleep(1.0)
        try:
            if dialogWrap.window_text() == 'iTunes' \
                    and dialog.Edit1.wait('ready').window_text() == 'Apple ID' \
                    and dialog.Edit2.wait('ready').window_text() == 'Password' \
                    and dialog.Button1.wait('exists').window_text() == '&Sign In':
                break
        except Exception as e:
            continue
    else:
        reportResult(error_code.REQ_LOGIN_INFO_ERR,"没有找到登录窗口")
        raise Exception("Failed to find login window in 15 iterations!")
    app.wait_cpu_usage_lower(cpu_usage)

    logger.info("Request login info from %s" % webAddress)

    login_result = False
    login_info_list = []
    for i in range(3):
        logger.info("do login index:%d" % i)
        for k in range(3):
            logger.info("request login info index:%d" % k)
            # 请求用户名和密码
            data_json = json.dumps({'task_id': taskId})
            url = webAddress + '/scriptLoginInfoRequest'
            responseData = requests.post(url, data_json)

            logger.info("scriptLoginInfoRequest status_code:%d" % responseData.status_code)
            if responseData.status_code == 200:
                break

            reportResult(error_code.REQ_LOGIN_INFO_ERR, "scriptLoginInfoRequest status_code:%d" % responseData.status_code)
            logger.info("sleep 5s and retry")
            time.sleep(5.0)
            continue
                

        # 判断内容是否和前面的一样，如果内容一样忽略登录
        new_login_info = True
        for item in login_info_list:
            if item == responseData.text:
                new_login_info = False
                break

        if new_login_info == True:
            login_info_list.append(responseData.text)
        else:
            logger.info("not request new login info, ignore login")
            time.sleep(15.0)
            continue

        jsonData = json.loads(responseData.text)
        appleId = jsonData["apple_id"]
        applePwd = jsonData["apple_pwd"]
        logger.info("request appleId:%s" % appleId)
        logger.info("request applePwd:%s" % applePwd)

        logger.info("Setting login dialog edit texts")

        appleid_Edit = dialog.Edit1
        appleid_Edit.wait('ready')
        appleid_Edit.click_input()
        appleid_Edit.type_keys(appleId)
        appleid_Edit.set_edit_text(appleId)
        time.sleep(3)

        pass_Edit = dialog.Edit2
        pass_Edit.wait('ready')
        pass_Edit.click_input()
        pass_Edit.type_keys(applePwd)
        pass_Edit.set_edit_text(applePwd)
        time.sleep(3)

        logger.info("Clicking login button!")
        loginButton = dialog.Button1
        loginButton.wait('ready')
        # click multiple times as pywinauto seems to have some bug
        loginButton.click()
        time.sleep(0.5)
        try:
            loginButton.click()
            time.sleep(0.5)
            loginButton.click_input()
        except:
            pass

        logger.info("Waiting login result...")
        time.sleep(10)
        debugText = debugTopWin()

        if "Sign In to the iTunes Store" in debugText:
            logger.info("Failed to trigger Login button!")
            reportResult(error_code.REQ_LOGIN_INFO_ERR,"没有点击登录按钮")
            raise Exception("Failed to trigger Login button!")
        elif app.top_window().window_text() == 'Verification Failed':
            logger.info("login Verification Failed")
            reportResult(error_code.REQ_LOGIN_INFO_ERR,"登录失败，请检查账号密码")
            # raise Exception("Verification Failed: %s" % app.top_window().Static2.window_text())
            cleanFailedDialog()
            time.sleep(15.0)
        else:
            logger.info("login success")
            login_result = True
            break

    if login_result == True:
        reportResult(error_code.REQ_LOGIN_INFO_SUCCESS, "")
    else:
        raise Exception("login Failed")
        # exit(error_code.REQ_LOGIN_INFO_ERR)


def tfaItunes():
    logger.info("Check 2FA auth...")
    need2FA = False
    for i in range(6):
        winText = debugTopWin()
        if "Enter the verification code sent to your other devices." in winText or "SPINNER" in winText:
            logger.info("need 2FA auth")
            need2FA = True
            dialog = app.top_window()
            dialogWrap = dialog.wait('ready')
            break
        else:
            logger.info("check 2FA auth sleep 3s")
            time.sleep(5.0)

    app.wait_cpu_usage_lower(cpu_usage)

    if need2FA == True:
        logger.info("need 2FA auth")
    else:
        logger.info("not need 2FA auth")

    if need2FA == False:
        return

    login_result = False
    login_info_list = []
    for i in range(3):
        logger.info("do 2FA login index:%d" % i)
        for k in range(12):
            time.sleep(5.0)
            logger.info("Start request 2FA from web index:%d" % k)

            data_json = json.dumps({'task_id': taskId})
            url = webAddress + '/script2FARequest'
            responseData = requests.post(url, data_json)

            logger.info("script2FARequest result:%d " % responseData.status_code)
            if responseData.status_code != 200:
                continue

            jsonData = json.loads(responseData.text)
            twoFACode = jsonData["two_fa_code"]
            logger.info("web 2FA is:%s" % twoFACode)

            if len(twoFACode) == 6:
                logger.info("request web 2FA success")
                break

            logger.info("not read 2FA from web, sleep 5s ,2FA len:%d" % len(twoFACode))
        else:
            reportResult(error_code.REQ_2FA_INFO_ERR, "not read 2FA in 60s")
            raise Exception("not read 2FA in 60s")

        # 判断内容是否和前面的一样，如果内容一样忽略登录
        new_login_info = True
        for item in login_info_list:
            if item == responseData.text:
                new_login_info = False
                break

        if new_login_info == True:
            login_info_list.append(responseData.text)
        else:
            logger.info("not request new 2FA info, ignore login 2FA")
            time.sleep(15.0)
            continue

        twoFA1 = twoFACode[0]
        twoFA2 = twoFACode[1]
        twoFA3 = twoFACode[2]
        twoFA4 = twoFACode[3]
        twoFA5 = twoFACode[4]
        twoFA6 = twoFACode[5]

        logger.info("Setting 2FA dialog edit texts")

        twoFA_Edit1 = dialog.Edit1
        twoFA_Edit1.wait('ready')
        twoFA_Edit1.click_input()
        twoFA_Edit1.type_keys(twoFA1)
        twoFA_Edit1.set_edit_text(twoFA1)
        time.sleep(1)

        twoFA_Edit2 = dialog.Edit2
        twoFA_Edit2.wait('ready')
        twoFA_Edit2.click_input()
        twoFA_Edit2.type_keys(twoFA2)
        twoFA_Edit2.set_edit_text(twoFA2)
        time.sleep(1)

        twoFA_Edit3 = dialog.Edit3
        twoFA_Edit3.wait('ready')
        twoFA_Edit3.click_input()
        twoFA_Edit3.type_keys(twoFA3)
        twoFA_Edit3.set_edit_text(twoFA3)
        time.sleep(1)

        twoFA_Edit4 = dialog.Edit4
        twoFA_Edit4.wait('ready')
        twoFA_Edit4.click_input()
        twoFA_Edit4.type_keys(twoFA4)
        twoFA_Edit4.set_edit_text(twoFA4)
        time.sleep(1)

        twoFA_Edit5 = dialog.Edit5
        twoFA_Edit5.wait('ready')
        twoFA_Edit5.click_input()
        twoFA_Edit5.type_keys(twoFA5)
        twoFA_Edit5.set_edit_text(twoFA5)
        time.sleep(1)

        twoFA_Edit6 = dialog.Edit6
        twoFA_Edit6.wait('ready')
        twoFA_Edit6.click_input()
        twoFA_Edit6.type_keys(twoFA6)
        twoFA_Edit6.set_edit_text(twoFA6)
        time.sleep(1)

        logger.info("Clicking 2FA button!")
        loginButton = dialog.Button1
        loginButton.wait('ready')
        # click multiple times as pywinauto seems to have some bug
        loginButton.click()
        time.sleep(0.5)
        try:
            loginButton.click()
            time.sleep(0.5)
            loginButton.click_input()
        except:
            pass

        logger.info("Waiting 2FA result...")
        time.sleep(10)
        debugTopWin()

        if app.top_window().handle == dialogWrap.handle:
            logger.info("Failed to trigger 2FA button!")
            reportResult(error_code.REQ_2FA_INFO_ERR,"没有点击2次认证按钮")
            raise Exception("Failed to trigger 2FA button!")
        elif app.top_window().window_text() == 'Verification Failed':
            logger.info("login 2fa Verification Failed")
            reportResult(error_code.REQ_2FA_INFO_ERR, "二次验证码错误，请重新输入")
            # raise Exception("Verification Failed: %s" % app.top_window().Static2.window_text())
            cleanFailedDialog()
        else:
            logger.info("login 2fa success")
            login_result = True
            break

    if login_result == True:
        reportResult(error_code.REQ_2FA_INFO_SUCCESS, "")
    else:
        raise Exception("login 2fa Failed")
        # exit(error_code.REQ_2FA_INFO_ERR)


def initITunes():
    subprocess.call('taskkill /f /im APSDaemon*', shell=True)
    subprocess.call('taskkill /f /im iTunes*', shell=True)

    global app
    app = Application().start(r"C:\Program Files\iTunes\iTunes.exe")
    app.wait_cpu_usage_lower(cpu_usage)
    time.sleep(8)

    # Click all first-time dialogs (like License Agreements, missing audios)
    cleanAllDialog()

    # Calm down a bit before main window operations
    app.wait_cpu_usage_lower(cpu_usage)
    debugTopWin()

    cleanWelcomeDialog()

    # TODO:Already logged in!需要处理
    loginItunes()
    time.sleep(5)
    tfaItunes()
    
    app.wait_cpu_usage_lower(cpu_usage)

    # Finish & Cleanup
    logger.info("Waiting all dialogs to finish")
    time.sleep(5)
    cleanAllDialog()
    reportResult(error_code.REQ_LOGIN_SUCCESS, "")



for init_i in range(3):
    if len(webAddress) == 0:
        logger.fatal("webAddress is empty, stop login task")
        exit(1)

    if len(taskId) == 0:
        logger.fatal("taskId is empty, stop login task")
        exit(2)

    try:
        initITunes()
        break
    except Exception as e:
        logger.info("Init iTunes %d: Failed with %s" % (init_i, e))
        import traceback; traceback.print_exc()
        time.sleep(8)

logger.info("Init iTunes Successfully!")