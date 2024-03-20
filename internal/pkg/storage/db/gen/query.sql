-- name: GetModels :many
select id,
       name,
       description,
       count(1) over () as count
from models
where name like $1
order by created_at desc
limit $2 offset $3;

-- name: GetModel :one
select id,
       name,
       description
from models
where id = $1;

-- name: GetModelProblem :one
select p.id,
       p.name,
       p.description,
       array_agg(me.id)          as metric_ids,
       array_agg(me.name)        as metric_names,
       array_agg(me.description) as metric_descriptions
from models m
         join problems p on m.problem_id = p.id
         join problem_metrics pm on p.id = pm.problem_id
         join metrics me on pm.metric_id = me.id
where m.id = $1
group by (p.id, p.name, p.description);

-- name: GetModelHyperparameters :many
select h.id,
       h.name,
       h.description,
       h.type,
       h.default_value
from models m
         join hyperparameters h on m.id = h.model_id
where m.id = $1;