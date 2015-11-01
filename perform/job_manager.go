package perform

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"
)

func RunDeployJobs(do *definitions.Do) error {
	for _, job := range do.Package.Jobs {
		logger.Printf("Executing Job Named =>\t\t%s\n", job.JobName)
		var err error
		switch {
		// Util jobs
		case job.Job.Account != nil:
			logger.Infof("\tType =>\t\t\tAccount\n")
			job.JobResult, err = SetAccountJob(job.Job.Account, do)
		case job.Job.Set != nil:
			logger.Infof("\tType =>\t\t\tSet\n")
			job.JobResult, err = SetValJob(job.Job.Set, do)
		// Transaction jobs
		case job.Job.Send != nil:
			logger.Infof("\tType =>\t\t\tSend\n")
			job.JobResult, err = SendJob(job.Job.Send, do)
		case job.Job.RegisterName != nil:
			logger.Infof("\tType =>\t\t\tRegisterName\n")
			job.JobResult, err = RegisterNameJob(job.Job.RegisterName, do)
		case job.Job.Permission != nil:
			logger.Infof("\tType =>\t\t\tPermission\n")
			job.JobResult, err = PermissionJob(job.Job.Permission, do)
		case job.Job.Bond != nil:
			logger.Infof("\tType =>\t\t\tBond\n")
			job.JobResult, err = BondJob(job.Job.Bond, do)
		case job.Job.Unbond != nil:
			logger.Infof("\tType =>\t\t\tUnbond\n")
			job.JobResult, err = UnbondJob(job.Job.Unbond, do)
		case job.Job.Rebond != nil:
			logger.Infof("\tType =>\t\t\tRebond\n")
			job.JobResult, err = RebondJob(job.Job.Rebond, do)
		// Contracts jobs
		case job.Job.Deploy != nil:
			logger.Infof("\tType =>\t\t\tDeploy\n")
			job.JobResult, err = DeployJob(job.Job.Deploy, do)
		case job.Job.PackageDeploy != nil:
			logger.Infof("\tType =>\t\t\tPackageDeploy\n")
			job.JobResult, err = PackageDeployJob(job.Job.PackageDeploy, do)
		case job.Job.Call != nil:
			logger.Infof("\tType =>\t\t\tCall\n")
			job.JobResult, err = CallJob(job.Job.Call, do)
		// State jobs
		case job.Job.RestoreState != nil:
			logger.Infof("\tType =>\t\t\tRestoreState\n")
			job.JobResult, err = RestoreStateJob(job.Job.RestoreState, do)
		case job.Job.DumpState != nil:
			logger.Infof("\tType =>\t\t\tDumpState\n")
			job.JobResult, err = DumpStateJob(job.Job.DumpState, do)
		}

		if err != nil {
			return err
		}

		if err = util.WriteJobResult(job.JobName, job.JobResult); err != nil {
			return err
		}
	}

	return nil
}

func RunTestJobs(do *definitions.Do) error {
	for _, job := range do.Package.Jobs {
		logger.Printf("Executing Job Named =>\t\t%s\n", job.JobName)
		var err error
		switch {
		case job.Job.Query != nil:
			logger.Infof("\tType =>\t\t\tQuery\n")
			job.JobResult, err = QueryJob(job.Job.Query)
		case job.Job.GetNameEntry != nil:
			logger.Infof("\tType =>\t\t\tGetNameEntry\n")
			job.JobResult, err = GetNameEntryJob(job.Job.GetNameEntry)
		case job.Job.Assert != nil:
			logger.Infof("\tType =>\t\t\tAssert\n")
			job.JobResult, err = AssertJob(job.Job.Assert)
		}

		if err != nil {
			return err
		}

		if err = util.WriteJobResult(job.JobName, job.JobResult); err != nil {
			return err
		}
	}

	return nil
}
