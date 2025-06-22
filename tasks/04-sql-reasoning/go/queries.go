// tasks/04‑sql‑reasoning/go/queries.go
package queries

// Task A
const SQLA = `
	SELECT
		p.campaign_id,
		SUM(p.amount_thb) AS total_thb,
		ROUND(SUM(p.amount_thb) * 1.0 / c.target_thb, 4) AS pct_of_target
	FROM
		pledge AS p
		LEFT JOIN campaign AS c ON p.campaign_id = c.id
	GROUP BY
		p.campaign_id
	ORDER BY
		pct_of_target DESC,
		campaign_id ASC
`

// Task B
const SQLB = `
	WITH group_pledges AS (
		SELECT 'global' AS scope, amount_thb
		FROM pledge 
		UNION ALL
		SELECT 'thailand' AS scope, amount_thb
		FROM pledge p JOIN donor d ON d.id = p.donor_id
		WHERE country = 'Thailand'
	), nearest_rank AS (
		SELECT scope, amount_thb, ROW_NUMBER() OVER (PARTITION BY scope ORDER BY amount_thb) AS rn,
			COUNT(*) OVER(PARTITION BY scope) AS total
		FROM group_pledges
	), p90 AS (
		SELECT scope, amount_thb AS p90_thb, rn, total, CAST((0.9 * total + 0.9999999) AS INT) AS p90_rank
		FROM nearest_rank
	)SELECT scope, p90_thb
	FROM p90
	WHERE rn = p90_rank
	ORDER BY CASE WHEN scope = 'global' THEN 1 ELSE 2 END
`

var Indexes = []string{} // skipped
