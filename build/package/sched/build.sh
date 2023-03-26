#!/usr/bin/env bash
cd ${ROOT_DIR}
make go.build.linux_amd64.sched
cp ${OUTPUT_DIR}/platforms/${IMAGE_PLAT}/mqtt_sender "${DST_DIR}"
cp -rv ${CFG_DIR}/sched.yaml "$DST_DIR"