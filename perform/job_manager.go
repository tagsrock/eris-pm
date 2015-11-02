package perform

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"
)

func RunDeployJobs(do *definitions.Do) error {
	for _, job := range do.Package.Jobs {
		var err error
		switch {

		// Util jobs
		case job.Job.Account != nil:
			announce(job.JobName, "Account")
			job.JobResult, err = SetAccountJob(job.Job.Account, do)
		case job.Job.Set != nil:
			announce(job.JobName, "Set")
			job.JobResult, err = SetValJob(job.Job.Set, do)

		// Transaction jobs
		case job.Job.Send != nil:
			announce(job.JobName, "Sent")
			job.JobResult, err = SendJob(job.Job.Send, do)
		case job.Job.RegisterName != nil:
			announce(job.JobName, "RegisterName")
			job.JobResult, err = RegisterNameJob(job.Job.RegisterName, do)
		case job.Job.Permission != nil:
			announce(job.JobName, "Permission")
			job.JobResult, err = PermissionJob(job.Job.Permission, do)
		case job.Job.Bond != nil:
			announce(job.JobName, "Bond")
			job.JobResult, err = BondJob(job.Job.Bond, do)
		case job.Job.Unbond != nil:
			announce(job.JobName, "Unbond")
			job.JobResult, err = UnbondJob(job.Job.Unbond, do)
		case job.Job.Rebond != nil:
			announce(job.JobName, "Rebond")
			job.JobResult, err = RebondJob(job.Job.Rebond, do)

		// Contracts jobs
		case job.Job.Deploy != nil:
			announce(job.JobName, "Deploy")
			job.JobResult, err = DeployJob(job.Job.Deploy, do)
		case job.Job.PackageDeploy != nil:
			announce(job.JobName, "PackageDeploy")
			job.JobResult, err = PackageDeployJob(job.Job.PackageDeploy, do)
		case job.Job.Call != nil:
			announce(job.JobName, "Call")
			job.JobResult, err = CallJob(job.Job.Call, do)

		// State jobs
		case job.Job.RestoreState != nil:
			announce(job.JobName, "RestoreState")
			job.JobResult, err = RestoreStateJob(job.Job.RestoreState, do)
		case job.Job.DumpState != nil:
			announce(job.JobName, "DumpState")
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
		var err error
		switch {
		case job.Job.QueryAccount != nil:
			announce(job.JobName, "QueryAccount")
			job.JobResult, err = QueryAccountJob(job.Job.QueryAccount, do)
		case job.Job.QueryContract != nil:
			announce(job.JobName, "QueryContract")
			job.JobResult, err = QueryContractJob(job.Job.QueryContract, do)
		case job.Job.QueryName != nil:
			announce(job.JobName, "QueryName")
			job.JobResult, err = QueryNameJob(job.Job.QueryName, do)
		case job.Job.QueryVals != nil:
			announce(job.JobName, "QueryVals")
			job.JobResult, err = QueryValsJob(job.Job.QueryVals, do)
		case job.Job.Assert != nil:
			announce(job.JobName, "Assert")
			job.JobResult, err = AssertJob(job.Job.Assert, do)
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

func announce(job, typ string) {
	logger.Printf("Executing Job Named =>\t\t%s\n", job)
	logger.Infof("\tType =>\t\t\t%s\n", typ)
}
