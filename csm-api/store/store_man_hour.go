package store

import (
	"context"
	"csm-api/entity"
	"fmt"
)

// func: 프로젝트에 설정된 공수 조회
// @param
// - jno: 프로젝트pk
func (r *Repository) GetManHourList(ctx context.Context, db Queryer, jno int64) (*entity.ManHours, error) {
	manHours := entity.ManHours{}

	query := `
			SELECT
			    MHNO,
			    WORK_HOUR,
			    MAN_HOUR,
			    JNO
			FROM 
			    IRIS_MAN_HOUR MH
			WHERE
				MH.JNO = :1
			ORDER BY
			    WORK_HOUR ASC
		`

	if err := db.SelectContext(ctx, &manHours, query, jno); err != nil {
		return nil, fmt.Errorf("GetManHourList err:%v", err)
	}

	return &manHours, nil
}

// func: 공수 수정 및 추가
// @param
// - manHour: 공수 정보
func (r *Repository) MergeManHour(ctx context.Context, tx Execer, manHour entity.ManHour) (err error) {
	query := `
		MERGE INTO SAFE.IRIS_MAN_HOUR J1
		USING (
			SELECT 
				:1 AS MHNO,
				:2 AS WORK_HOUR,
				:3 AS MAN_HOUR,
				:4 AS JNO, 
				:5 AS UNO,	
				:6 AS USER_NAME
			FROM DUAL
		) J2
		ON (
			J1.MHNO = J2.MHNO
		) WHEN MATCHED THEN
			UPDATE SET
				J1.WORK_HOUR = J2.WORK_HOUR,
				J1.MAN_HOUR = J2.MAN_HOUR,
				J1.JNO = J2.JNO,
				J1.MOD_UNO = J2.UNO,	
				J1.MOD_USER = J2.USER_NAME,
				J1.MOD_DATE = SYSDATE
		WHEN NOT MATCHED THEN
			INSERT ( MHNO, WORK_HOUR, MAN_HOUR, JNO, REG_UNO, REG_USER, REG_DATE )
			VALUES (
				SEQ_IRIS_MAN_HOUR.NEXTVAL,
				J2.WORK_HOUR,
				J2.MAN_HOUR,
				J2.JNO,
				J2.UNO,	
				J2.USER_NAME,
				SYSDATE
			)
		`
	if _, err = tx.ExecContext(ctx, query, manHour.Mhno, manHour.WorkHour, manHour.ManHour, manHour.Jno, manHour.RegUno, manHour.RegUser); err != nil {
		//TODO: 에러 아카이브
		return fmt.Errorf("MargeManHour err: %w", err)

	}

	return
}
