UPDATE
    events
SET
    title = $2,
    date = $3,
    duration = $4,
    description = $5,
    userid = $6,
    notify = $7
WHERE
    id = $1