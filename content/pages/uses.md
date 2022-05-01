---
title: "Uses"
subtitle: "Behind the curtain"
herotext: "These are the tools that power my workflows"
summary: "These are the tools that power my workflows"
---

## Computers

- AMD Ryzen 3950X / 80GB RAM / 9TB Disks - Running Windows today, but that changes roughly weekly. I'm a _Professional_ at installing new operating systems and getting up and running with my code & tools within an hour. In the last two weeks this computer has had Ubuntu 22.04, Fedora Silverblue, Pop_OS! 22.04, and Windows 11. (I was chasing down a networking problem that ended up being a bad switch and two bad ethernet cables)

- Surface Laptop 3 (Business Edition) - Since I'm not traveling lately, this doesn't see much use beyond the occasional couch coding session. It's running Windows right now, but I like running Pop_OS! on it too. It's really a great laptop. The size is perfect and I like the aspect ratio of the screen. 32GB of memory and a 1TB NVMe drive are plenty for anything I do on it.

- Monitors - LG 27" 4k and Phillips 34" curved. Mounted one on top of the other on the desk - see below.

## Software

- VS Code - It's really impressive how capable this editor is and how it manages to feel lightweight while having features far beyond your average text editor.
- [Multipass](https://multipass.run) - I use Multipass instead of WSL on the Windows desktop. Multipass on Windows uses Hyper-V, and allows you to add multiple networks to the VM, so I've added a bridged adapter in addition to the private network Multipass assigns. The bridged adapter exposes the VM on my LAN so it's available by SSH from other computers. Doing this has allowed me to forego my usual Linux workstation setup. I gave the primary VM a lot of RAM, CPU and disk so it is a very capable development environment. I'm surprised I don't hear or see more about Multipass from others, this is an excellent solution for developers.
- Microsoft Edge - I've grown quite fond of Microsoft Edge and I use it as my primary browser on everything from my iPhone/iPad to my desktop/laptop.

## Accessories

- [Moonlander Keyboard](https://www.zsa.io/moonlander/) - `<3 <3 <3 <3` I was having a lot of RSI problems until I switched to the Moonlander and trackball. Now they're completely gone. Take care of yourself, kids.
- Logitech Wireless Trackball - the wireless is spotty, I think my magnetic desktop messes with the signal. When I relocate the trackball just a few inches away from where it was reception works again. It isn't drivers - the behavior is the same on all operating systems. It drives me batty, but I can't find a thumb ball trackball other than this one that I like. If you have suggestions hit me up on twitter @bketelsen.

## Video and Streaming

- Logitech Brio web cam - for the Teams/Zoom/etc meetings. It's pretty decent in Windows. Not amazing in macOS, OK in Linux.
- Rode Podcaster Pro sound board - great board, but it's missing some output options I'd like - like anything line-level. It's annoying taking headphone level output and plugging it into the ATEM. That's a crime against audio.
- Canon XA11 HD Camcorder - for the more important video work
- Blackmagic Atem Mini Extreme - small form factor, has most of the video effects of my Tricaster Mini HD which is relegated to a cupboard
- [NewBlue Titler Live Broadcast](https://newbluefx.com/products/titler-live/broadcast/) - bought it on a whim, but haven't invested enough time to become really good with it yet. It's really powerful and feeds titles into the ATEM perfectly
- Several various NDI to HDMI adapters so I can feed remote video (say from a mac or Linux box or even Titler Live) to the ATEM without screwing with cables or splitting HDMI cables
- Stream Deck & Stream Deck XL

## Desk

This is where I broke the bank. I have an [AltWork signature](https://altwork.com/products/altwork-signature-station) that I got for a great discount as a demo unit. It has some scratches and dings, but it works great and really takes the pressure off my damaged spinal column. A life-saver for a programmer/desk worker with spinal degradation. It comes with a _LOT_ of limitations - cables are required to be really long unless you're using a laptop (I'm not), and there's only so many cables you can jam through their cable management guides. That can be frustrating when all your accessories are wired and there are limitations to how long you can transmit video & data over a USB-C cable. Instead I run HDMI and DisplayPort cables. It's a little awkward, but it works. The desk surface is magnetic so I keep the keyboard and trackball from moving while I'm reclined by attaching rare earth magnets to them. I only forgot about my coffee cup when I was reclining once. That's all it took :)

## NAS

Synology DS-1817+ - it's a little underpowered to do much more than be a NAS and run a docker container or two. If I had it to do again, I'd likely build something more powerful as a custom build. But it's been rock solid for years, and it has 28TB of storage that I've filled about half way with packrat files and backups.

## Network Gear

Unifi Gear everywhere. Dream Station Pro, 24 Port switch, several (5) APs, IP cameras and more I've probably forgotten. It works mostly, but I get random strange things happening which tempts me to rewire the old CAT5 in the house and start the network configuration from scratch. I'm running an isolated guest network and an internal network. I could probably do a lot more with VLANs if I took the trouble. I really want to upgrade this to faster than 1GB speeds. That gets complicated quickly though. Maybe just a 10G switch in my office with all the important stuff would suffice for a while.

Last update 1 May 2022
