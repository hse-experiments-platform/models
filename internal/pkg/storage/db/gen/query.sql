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
       array_agg(me.metric_name) as metric_names
from models m
         join problems p on m.problem_id = p.id
         join metrics me on m.id = me.model_id
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
       l.launch_status,
       m.id                as model_id,
       m.name              as model_name,
       p.name              as problem_name,
       tm.train_dataset_id as training_dataset_id,
       d.name              as training_dataset_name,
       tm.created_at,
       tm.launch_id,
       count(1) over ()    as count
from trained_models tm
         join models m on tm.base_model_id = m.id
         join problems p on m.problem_id = p.id
         join datasets d on d.id = tm.train_dataset_id
         join launches l on tm.launch_id = l.id
where tm.name like $1
  and (tm.base_model_id = sqlc.arg(model_id) or sqlc.arg(model_id) = 0)
order by created_at desc
limit $2 offset $3;

-- name: GetTrainedModel :one
select tm.id,
       tm.name,
       tm.description,
       l.launch_status,
       m.id          as model_id,
       m.name        as model_name,
       p.id          as problem_id,
       p.description as problem_description,
       p.name        as problem_name,
       tm.train_dataset_id,
       d.name        as training_dataset_name,
       tm.created_at,
       tm.launch_id,
       tm.target_column
from trained_models tm
         join models m on tm.base_model_id = m.id
         join problems p on m.problem_id = p.id
         join datasets d on d.id = tm.train_dataset_id
         join launches l on tm.launch_id = l.id
where tm.id = $1;

-- name: GetAllModels :many
select m.id,
       m.name,
       m.description,
       p.name                              as problem_name,
       m.class_name,
       array_remove(h.name, null)          as hyperparameter_names,
       array_remove(h.description, null)   as hyperparameter_descriptions,
       array_remove(h.type, null)          as hyperparameter_types,
       array_remove(h.default_value, null) as hyperparameter_default_values,
       array_remove(m2.metric_name, null)  as metric_names
from models m
         join problems p on m.problem_id = p.id
         cross join lateral (select array_agg(h.name)          as name,
                                    array_agg(h.description)   as description,
                                    array_agg(h.type)          as "type",
                                    array_agg(h.default_value) as default_value
                             from hyperparameters h
                             where m.id = h.model_id) as h
         cross join lateral (select array_agg(m2.metric_name) as metric_name
                             from metrics m2
                             where m.id = m2.model_id) as m2;

-- name: CreateModel :one
insert into models (name, description, problem_id, class_name)
values (sqlc.arg(name), sqlc.arg(description), (select id from problems where name = sqlc.arg(problem)), sqlc.arg(class_name))
returning id;

-- name: CreateHyperparameters :exec
insert into hyperparameters (name, description, type, default_value, model_id)
select unnest(sqlc.arg(names)::text[]),
       unnest(sqlc.arg(descriptions)::text[]),
       unnest(sqlc.arg(types)::text[]),
       unnest(sqlc.arg(default_values)::text[]),
       $5;

-- name: CreateModelMetrics :exec
insert into metrics (metric_name, model_id)
select unnest(sqlc.arg(metric_names)::text[]), $2;

-- name: GetTrainedModelRunID :one
select output
from trained_models tm
         join launches l on tm.launch_id = l.id
where tm.id = sqlc.arg(trained_model_id);
