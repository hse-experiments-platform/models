-- name: GetModels :many
select id,
       name,
       description,
       count(1) over () as count
from models
where name like $1
  and (problem_id = $4 or $4 = 0)
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

-- name: GetProblems :many
select id,
       name,
       description,
       count(1) over () as count
from problems
where name like $1
order by created_at desc
limit $2 offset $3;

-- name: GetTrainedModels :many
select tm.id,
       tm.name,
       tm.description,
       tm.model_training_status,
       m.id             as model_id,
       m.name           as model_name,
       p.name           as problem_name,
       tm.training_dataset_id as training_dataset_id,
       d.name           as training_dataset_name,
       tm.created_at,
       tm.launch_id,
       count(1) over () as count
from trained_models tm
         join models m on tm.model_id = m.id
         join problems p on m.problem_id = p.id
         join datasets d on d.id = tm.training_dataset_id
where tm.name like $1
  and (tm.model_id = sqlc.arg(model_id) or sqlc.arg(model_id) = 0)
order by created_at desc
limit $2 offset $3;

-- name: GetTrainedModel :one
select tm.id,
       tm.name,
       tm.description,
       tm.model_training_status,
       m.id          as model_id,
       m.name        as model_name,
       p.id          as problem_id,
       p.description as problem_description,
       p.name        as problem_name,
       tm.training_dataset_id,
       d.name        as training_dataset_name,
       tm.created_at,
       tm.launch_id,
       tm.target_column
from trained_models tm
         join models m on tm.model_id = m.id
         join problems p on m.problem_id = p.id
         join datasets d on d.id = tm.training_dataset_id
where tm.id = $1;
