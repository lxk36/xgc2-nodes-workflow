# xgc2-nodes-workflow

Product-neutral control nodes for `xgc2-orchestration-core`.

The `actioncall` descriptor marks a Workflow occurrence whose concrete input,
trigger, scope, and result schemas are frozen by `CallAction`. The public core,
not the extension executor, creates and joins the durable child Run. This keeps
child lineage, exact Action pins, retries, and ownership closure inside the
orchestration kernel while leaving the control node independently publishable.

The `signalwait` node gives long-running workflows an explicit durable hold
point. Event ingress must match its exact subject and condition digest before
the kernel resumes the same invocation. A workflow can therefore start managed
processes, remain waiting while an experiment is active, and run exact stop
nodes after a user or policy signal without a hidden product scheduler.
