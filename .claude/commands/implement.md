Implement the tasks defined.

This is the fourth step in the Spec-Driven Development lifecycle.

Given the context provided as an argument, do this:

1. Read the `.spec-kit/specs/[###-feature-name]/tasks.md`

2. Execute the task in order.

    - If the task has [P] symbol, execute the task in parallel

3. Validate if the task is well executed, perform any test, you consider necessary to check if is done.

    - If is incorrect retry the task.
    - Run more tests, to validate the retry.

4. Check for the next incomplete task.

    - the task must be executed in order.

5. Repeat the process until all the task are completed.
