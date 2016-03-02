package util

import (
	"fmt"
	"regexp"

	"github.com/eris-ltd/eris-pm/definitions"

	log "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func MintChainErrorHandler(do *definitions.Do, err error) (string, error) {
	log.WithFields(log.Fields{
		"defAddr":  do.Package.Account,
		"chainID":  do.ChainID,
		"chainURL": do.Chain,
		"rawErr":   err,
	}).Error("")

	return "", fmt.Errorf(`
There has been an error talking to your eris chain.

%v

Debugging this error is tricky, but don't worry the marmot recovery checklist is...
  * is the %s account right?
  * is the account you want to use in your keys service: eris keys ls ?
  * is the account you want to use in your genesis.json: eris chains cat %s genesis ?
  * is your chain making blocks: eris chains logs -f %s ?
  * do you have permissions to do what you're trying to do on the chain?
`, err, do.Package.Account, do.ChainID, do.ChainID)
}

func KeysErrorHandler(do *definitions.Do, err error) (string, error) {
	log.WithFields(log.Fields{
		"defAddr": do.Package.Account,
	}).Error("")

	r := regexp.MustCompile(fmt.Sprintf("open /home/eris/.eris/keys/data/%s/%s: no such file or directory", do.Package.Account, do.Package.Account))
	if r.MatchString(fmt.Sprintf("%v", err)) {
		return "", fmt.Errorf(`
Unfortunately the marmots could not find the key you are trying to use in the keys service.

There are two ways to fix this.
  1. Import your keys from your host: eris keys import %s
  2. Import your keys from your chain:

eris chains exec %s "mintkey eris chains/%s/priv_validator.json" && \
eris services exec keys "chown eris:eris -R /home/eris"

Now, run  eris keys ls  to check that the keys are available. If they are not there
then change the account. Once you have verified that the keys for account

%s

are in the keys service, then rerun me.
`, do.Package.Account, do.ChainID, do.ChainID, do.Package.Account)
	}

	return "", fmt.Errorf(`
There has been an error talking to your eris keys service.

%v

Debugging this error is tricky, but don't worry the marmot recovery checklist is...
  * is your %s account right?
  * is the key for %s in your keys service: eris keys ls ?
`, err, do.Package.Account, do.Package.Account)
}

func ABIErrorHandler(do *definitions.Do, err error, call *definitions.Call, query *definitions.QueryContract) (string, error) {
	switch {
	case call != nil:
		log.WithFields(log.Fields{
			"data":   call.Data,
			"abi":    call.ABI,
			"dest":   call.Destination,
			"rawErr": err,
		}).Error("ABI Error")
	case query != nil:
		log.WithFields(log.Fields{
			"data":   query.Data,
			"abi":    query.ABI,
			"dest":   query.Destination,
			"rawErr": err,
		}).Error("ABI Error")
	}

	return "", fmt.Errorf(`
There has been an error in finding your ABI. ABI's are "Application Binary Interface"
and they are what let us know how to talk to smart contracts. These little json files
can be read by a variety of things which need to talk to smart contracts so they are
quite necessary to be able to find.

Usually this error is the result of a bad deploy event. Sometimes the marmots are
not told that there has been an error when there really was and do not stop a running
of jobs. The ABIs are saved after the deploy events. So if there was a glitch in the
matrix, we apologize in advance.

The marmot recovery checklist is...
  * ensure that your contracts successfully deployed
  * if you used imports you may need to correct the instance variable
  * if you have more than one contract in a single file you may need to correct the instance variable
`)
}
