import winreg
import os

def set_computer_name(new_name):
    try:
        # 打开注册表项
        key = winreg.OpenKey(winreg.HKEY_LOCAL_MACHINE, r"SYSTEM\CurrentControlSet\Control\ComputerName\ComputerName", 0, winreg.KEY_SET_VALUE)
        
        # 修改注册表值
        winreg.SetValueEx(key, "ComputerName", 0, winreg.REG_SZ, new_name)
        
        # 关闭注册表项
        winreg.CloseKey(key)
        
        # 通知系统名称已更改
        os.system(f'WMIC computersystem where name="%COMPUTERNAME%" call rename "{new_name}"')
        
        print(f"系统名称已成功修改为: {new_name}")
    except Exception as e:
        print(f"修改系统名称时出错: {e}")

# 示例调用
new_name = "NewComputerName"
set_computer_name(new_name)
