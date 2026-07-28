[Features]
- `Service` now captures the binding metadata Cloud Foundry sends but the decoder previously discarded: `SyslogDrainURL`, `InstanceGUID`, `InstanceName`, `BindingGUID` and `BindingName`. `SyslogDrainURL` is what `cf cups -l` sets on a user-provided instance, and `InstanceName` recovers the instance name for a named binding, where `Name` holds the binding name instead. (#25)
