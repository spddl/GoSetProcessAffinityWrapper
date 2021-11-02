# GoSetProcessAffinity as a wrapper with GUI

basically this is a configuration tool for the IFEO settings. Settings that go beyond the IFEO are done via a wrapper.
The wrapper sets itself as a "debugger" to be started before the actual program / game and then to start the exe with the requirements.

It is also possible to start scripts at the start as admin or system user or at the end via "MonitorProcess".

To start the process the IFEO entry "Debugger" is cleared before and set again after startup for this process, this may not be necessary for many applications so they can be started directly with the flag "DEBUG_ONLY_THIS_PROCESS" using the option "PassThrough" (not recommended for games with anticheat)

- CPU priority
    - idle to High is set via IFEO
    - Realtime is set via the wrapper (requires admin rights)
- Memory/Page priority
    - Idle to Normal is set via IFEO
- IO priority
    - Very Low to Normal is set via IFEO
    - High is set via the wrapper (requires admin rights)
- Boost is set via the wrapper (requires admin rights)

![screenshot](https://i.imgur.com/wmj6evA.png)