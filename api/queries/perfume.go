package queries

const SelectPerfumesForUser string = `
	SELECT 
		perfume.id, 
		perfume.name,
		perfume.house, 
		perfume.url, 
		perfume.is_empty,
		perfume.description,
		array_agg(perfume_metric.id) AS notes FROM perfume_metric
	JOIN perfume ON perfume.id = perfume_metric.perfume_id
	WHERE perfume.user_id = $1
	GROUP BY perfume.id
`

const CreatePerfumeForUser string = `
	INSERT INTO perfume (user_id, name, house, url, description, is_empty) 
	VALUES
	($1, $2, $3, $4, $5, $6)
	RETURNING *;
`

const CreatePerfumeWearForUser = `
	INSERT INTO wear (user_id, perfume_id)
	VALUES
	($1, $2)
`

const UpdatePerfumeForUser string = `
	UPDATE perfume SET 
	name = $2, house = $3, url = $4, description = $5, is_empty = $6
	WHERE id = $1
	RETURNING *;
`

const GetMetricsForPerfume string = `
	SELECT * FROM perfume_metric 
	WHERE perfume_id = $1
`

const DeleteMetricsForPerfume string = `
	DELETE FROM perfume_metric
	WHERE perfume_id = $1
`

const SelectPerfumeVectorsForUser string = `
	SELECT 
	DISTINCT ON (metrics.perfume_id)
	metrics.perfume_id, 
	array[t_val, g_val, f_val, m_val] AS vector, 
	extract(epoch from wear.created_at) epoch
		FROM (
			SELECT 
				perfume.id as perfume_id, 
				avg(metric.t) AS t_val,
				avg(metric.g) AS g_val,
				avg(metric.f) AS f_val,
				avg(metric.m) AS m_val FROM perfume_metric
			JOIN perfume ON perfume.id = perfume_metric.perfume_id
			JOIN metric ON metric.id = perfume_metric.note_id
			WHERE perfume.user_id = $1
			GROUP BY perfume.id
		) metrics
		LEFT JOIN wear on wear.perfume_id = metrics.perfume_id
	ORDER BY metrics.perfume_id, epoch desc
`

const SelectMoodVectorForUser string = `
	SELECT array[t_val, g_val, f_val, m_val] AS vector
		FROM (
			SELECT  
				avg(metric.t) AS t_val,
				avg(metric.g) AS g_val,
				avg(metric.f) AS f_val,
				avg(metric.m) AS m_val FROM metric
			WHERE metric.id IN (%s)
		)
`

const SelectFragranceVectors string = `
		SELECT * FROM metric 
		WHERE metric.type = 'perfume'
`
