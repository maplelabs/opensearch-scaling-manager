### This directory contains log files logged using the "logrus" module

#### How to use ?
1. Import the package as 'log "scaling_manager/logger"' in your module
2. Use the package alias "log" at different levels
   Ex: log.Info(msg), log.Error(msg), log.Warn(msg)
3. Also supports the type of log. Specify the type as first argument while calling the method.
   Ex: log.Info(log.ProvisionerInfo, msg)
   This prints the info type along with the msg passed.
   Refer the variables declared in logger module for the log types defined


Future scope:
1. Add a config file from which the log configuration can be read from.
2. Add handlers for Trace, Debug and Panic
