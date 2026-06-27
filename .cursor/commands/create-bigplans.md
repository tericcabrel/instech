# Create Big Plans

---
description: When the task is too big, create a big plan to break it down into smaller steps. This is a markdown file that will be used to track the progress of the task. The user will use another command to build each phase of the plan from the markdown file.
---

You are operating in **strict agent mode**.

Extract the following from the user input:

* Plan Path: the path of the cursor plan that will be exported as a .md file.

If either is missing:
→ STOP and ask the user for the missing information.

---

## Rules

* Do NOT write any code
* Do NOT modify files
* Only produce the markdown file from the plan

---

## Task

Using:

* Plan Path: extracted value

Generate a master plan (Markdown file) containing all the phases or steps; for each phase add a status that can be "pending" (the default) or "implemented"
The value will be updated when the phase implementation is done.

If the plan file list the steps, rename them "phases" in the new markdown file. Example: "Step 1" will become "Phase 1", "Step 2" will become "Phase 2", etc.

Add a Phase overview table at the top — all phases with pending status (update to implemented when done).
When necessary, add checklists within each phase for granular tasks.

The plan should be stored in the docs/bigplans/ folder.
File name should be the same as the Plan Path.