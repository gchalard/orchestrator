# Description
Kuberenetes native pipeline orchestrator.

# Components

- Controlplane
- Job
- [UI]
- [SCM Watcher]

## Controlplane
1. Spawn jobs based on workflows and triggers
2. Receive updates and logs from the jobs

### Spawning jobs
The spawning process is :
Generate an ephemeral unique credentials so traffic between controlplane and job is encrypted and can be read only by authorized entities.
The unique credentials will be used by the jobs to authenticated against the controlplane to:
1. Request the Job definition
2. Update its status at each Step

### Workflows
The goal is to be github actions syntax compatible
A workflow run is :
```json
{
    name: str,
    jobs: {
        job_id (str): Job
    }
}
```

### Jobs
The goal is to be github actions syntax compatible
A job is:
```json
{
    name: str,
    needs: list[Job] # dependencies between job to compute the tree of jobs,
    runs_on: image_uri, # The image to use in the kubernetes Job definition. It must be derived from the original image of this Job system else the job will fail.
    steps: list[Step]
}
```

### Steps
The goal is to be github actions syntax compatible
A step is:
```json
{
    name: str,
    run: str, # a bash command to run
    working-directory: str, # the directory to set as current directory before running the command
    env: map[str, str], # the environment variables to set before running the command
    shell: str # the shell to use at each step, defaults to bash. The relation between shells and step is 1:1. At each step we pop a new shell that will be closed at the end of the step so there is isolation between them.
}
```

## Job
1. Request the steps to run after authenticating
2. Run the steps
3. Update the status after each step
4. Update the status of the job after every step is done
5. Streams the steps logs to the controlplane 