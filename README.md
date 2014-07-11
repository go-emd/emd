Go-Emd - Golang - Embarassingly Distributed
==

## Description
This is a framework built for a clustered environment.  The Distribution amoung the cluster is made up of Workers (go routines that perform raw computing) and leaders (one per node) controls persistance and helps coordinate between nodes and the manager go routine giving metrics and the ability to control the distribution remotely via a web gui.

## Install and Use
1. export GOPATH=<path to root>/Go-Emd/emd:$GOPATH
2. Download the Go-Emd-Examples repo and begin playing.

## Releases
- v0.1: Only works in a single node environment.  Not heavily tested... But is the basis of what the rest of the framework will become.
- v0.2: Works minimally on multiple machines in a distribution.  Not heavily tested... But is the basis for external communication between nodes in the distribution.
- v0.3: Refactor of internal and external connectors between node leaders and workers.
- v0.4: Leader is capable to obtain status updates of each child worker and start/stop each worker.  Successfully tested on both windows and linux variants.
- v0.5: Rest endpoints are completed for each nodeleader, metrics, status's, starting and stopping of the distribution is capable.  Tested on Windows and Linux variants.
- v0.6: Remodel of emd seperating the emd library and imlementations of it.  Lots of misc. updates allowing for more optimized and smoother code reusability and maintainability.
- v0.7: Addition of gob encoding allowing backwards compatible communication externally.  Many optimization edits.  Added generic's to all connector interfaces.  Workers need to be updated since then.
