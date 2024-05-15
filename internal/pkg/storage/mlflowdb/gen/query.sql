-- name: GetMetrics :many
select key,
       value,
       step
from metrics
where run_uuid = sqlc.arg(run_id)
order by step, key;

