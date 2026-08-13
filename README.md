# xgc2-nodes-workflow

Product-neutral control nodes for `xgc2-orchestration-core`.

The `actioncall` descriptor marks a Workflow occurrence whose concrete input,
trigger, scope, and result schemas are frozen by `CallAction`. The public core,
not the extension executor, creates and joins the durable child Run. This keeps
child lineage, exact Action pins, retries, and ownership closure inside the
orchestration kernel while leaving the control node independently publishable.
