<![CDATA[# VMware Setup Guide

How to get the best experience running Ash Linux in VMware Workstation (Windows) or VMware Fusion (macOS).

## VM Settings

Before booting Ash, configure these settings in VMware:

| Setting | Value |
|---------|-------|
| **RAM** | 8 GB recommended (4 GB minimum) |
| **CPU Cores** | 4+ cores recommended |
| **Disk** | 30 GB+ (thin provisioned is fine) |
| **3D Acceleration** | ✅ Enable (VM Settings → Display) |
| **Video Memory** | 4 GB recommended |

## Required: VMX Display Fix

Ash uses **Hyprland**, a Wayland compositor. VMware's default Vulkan renderer causes tearing, black boxes, and freezes with Hyprland. You need to add two lines to your VM's `.vmx` file on the **host** machine.

### How to Edit the VMX File

1. **Shut down the VM** (not suspend — fully shut down)
2. **Find the `.vmx` file:**
   - **Windows**: Usually in `C:\Users\YourName\Virtual Machines\Ash\Ash.vmx`
   - **macOS**: Right-click the VM in Fusion → Show in Finder → right-click the `.vmxbundle` → Show Package Contents → find the `.vmx` file
3. **Open in a text editor** and add these lines at the end:

```
mks.enableVulkanRenderer = "FALSE"
svga.disableFIFO = "TRUE"
```

4. **Save** and **boot the VM**

### What These Settings Do

- `mks.enableVulkanRenderer = "FALSE"` — Disables VMware's Vulkan renderer, which doesn't work well with Wayland compositors. Falls back to the SVGA GPU, which is stable.
- `svga.disableFIFO = "TRUE"` — Prevents rendering artifacts in the framebuffer.

## Optional: Enable Clipboard Sharing

By default, VMware blocks clipboard sharing. To enable copy/paste between host and VM, add these lines to the same `.vmx` file:

```
isolation.tools.copy.disable = "FALSE"
isolation.tools.paste.disable = "FALSE"
isolation.tools.setGUIOptions.enable = "TRUE"
```

`open-vm-tools` is already installed inside the VM. After adding these lines:

```bash
# Inside the VM, restart the guest tools:
sudo systemctl restart open-vm-tools
```

## Troubleshooting

### Black Screen After Boot

The most common cause is the missing VMX display fix:

1. Shut down the VM
2. Add `mks.enableVulkanRenderer = "FALSE"` and `svga.disableFIFO = "TRUE"` to the `.vmx` file
3. Boot again

If it still doesn't work, check the Hyprland log:
```bash
# SSH into the VM or switch to TTY2 with Ctrl+Alt+F2
cat ~/.local/share/hyprland/hyprland.log | tail -50
```

### Screen Tearing / Visual Glitches

Ensure both VMX lines are present. If tearing persists, try adding:
```
svga.maxWidth = "1920"
svga.maxHeight = "1080"
```

### Clipboard Still Not Working

1. Verify all three `isolation.tools.*` lines are in the `.vmx`
2. Restart the VM tools: `sudo systemctl restart open-vm-tools`
3. Check that `vmtoolsd` is running: `systemctl status open-vm-tools`

### VM Runs Slowly

- Ensure 3D acceleration is enabled in VM Settings → Display
- Allocate at least 4 CPU cores
- Give the VM at least 8 GB RAM
- Close other VMs if running multiple

---

**Next:** [How Search Works →](how-search-works.md) | [Troubleshooting →](troubleshooting.md)
]]>
