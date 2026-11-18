package repositories

import (
	"context"
	"encoding/json"
	"strings"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/application"
	"payment-engine/pkg/errors"

	gopostgres "github.com/HADLakmal/go-postgres"
	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/container"
)

// PlanRepository reads and writes the subscription plan business data table.
type PlanRepository struct {
	sqlRunner
}

func (r *PlanRepository) Init(c container.Container) error {
	r.db = c.Resolve(application.ModuleSQL).(gopostgres.DatabaseReporter)
	return nil
}

// planColumns is the column order scanPlan reads and CreatePlan writes.
var planColumns = []string{columnParamKey, columnParamCode, columnParamKind, columnParamStatus, columnParamPayload}

func scanPlan(scan func(dest ...any) error) (events.SubscriptionPlanPayload, error) {
	var plan events.SubscriptionPlanPayload
	var key string
	var payload []byte
	if err := scan(&key, &plan.Code, &plan.Kind, &plan.Status, &payload); err != nil {
		return plan, errors.Wrap(err, "failed to scan plan")
	}
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &plan); err != nil {
			return plan, errors.Wrap(err, "failed to unmarshal plan payload")
		}
	}
	plan.ID = parseUUID(key)
	return plan, nil
}

func (r *PlanRepository) ListPlans(ctx context.Context, kind, status string) ([]events.SubscriptionPlanPayload, error) {
	where := []string{columnParamDeletedAt + ` IS NULL`}
	var args []any
	if kind != "" {
		args = append(args, kind)
		where = append(where, columnParamKind+` = $`+itoa(len(args)))
	}
	if status != "" {
		args = append(args, status)
		where = append(where, columnParamStatus+` = $`+itoa(len(args)))
	}
	query := `SELECT ` + columnList(planColumns) + ` FROM ` + domain.SubscriptionPlanTable + ` WHERE ` + strings.Join(where, ` AND `) +
		` ORDER BY ` + columnParamKind + ` DESC, (` + columnParamPayload + `->>'price')::numeric, ` + columnParamCode
	var plans []events.SubscriptionPlanPayload
	err := r.queryRows(ctx, query, args, func(rows gopostgres.RowsInterface) error {
		plan, err := scanPlan(rows.Scan)
		if err != nil {
			return errors.Wrap(err, "failed to read plan row")
		}
		plans = append(plans, plan)
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list plans")
	}
	return plans, nil
}

func (r *PlanRepository) getPlan(ctx context.Context, where string, arg string) (events.SubscriptionPlanPayload, error) {
	var plan events.SubscriptionPlanPayload
	var scanErr error
	found := false
	err := r.queryRows(ctx, `SELECT `+columnList(planColumns)+` FROM `+domain.SubscriptionPlanTable+` WHERE `+where+` AND `+columnParamDeletedAt+` IS NULL`,
		[]any{arg}, func(rows gopostgres.RowsInterface) error {
			plan, scanErr = scanPlan(rows.Scan)
			found = true
			return scanErr
		})
	if err != nil {
		return plan, errors.Wrap(err, "failed to get plan")
	}
	if !found {
		return plan, errors.NotFound("subscription plan")
	}
	return plan, nil
}

func (r *PlanRepository) GetPlanByCode(ctx context.Context, code string) (events.SubscriptionPlanPayload, error) {
	return r.getPlan(ctx, columnParamCode+` = $1`, code)
}

func (r *PlanRepository) GetPlanByID(ctx context.Context, id string) (events.SubscriptionPlanPayload, error) {
	return r.getPlan(ctx, columnParamKey+` = $1::uuid`, id)
}

func (r *PlanRepository) CreatePlan(ctx context.Context, plan events.SubscriptionPlanPayload) error {
	payload, err := toJSON(plan)
	if err != nil {
		return errors.Wrap(err, "failed to encode plan payload")
	}
	_, err = r.exec(ctx, insertQuery(domain.SubscriptionPlanTable, planColumns),
		plan.ID.String(), plan.Code, plan.Kind, plan.Status, payload)
	if err != nil {
		return errors.Wrap(err, "failed to create plan")
	}
	return nil
}

func (r *PlanRepository) UpdatePlan(ctx context.Context, plan events.SubscriptionPlanPayload) error {
	payload, err := toJSON(plan)
	if err != nil {
		return errors.Wrap(err, "failed to encode plan payload")
	}
	_, err = r.exec(ctx, `UPDATE `+domain.SubscriptionPlanTable+` SET `+setClause(2, columnParamKind, columnParamStatus, columnParamPayload)+
		` WHERE `+columnParamKey+` = $1::uuid`,
		plan.ID.String(), plan.Kind, plan.Status, payload)
	if err != nil {
		return errors.Wrap(err, "failed to update plan")
	}
	return nil
}
