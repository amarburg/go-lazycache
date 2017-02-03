
curl  -H 'Content-Type: application/json' -H  'Authorization: Bearer 2dbeedbabf7fa93c449f091de9fe3ecda795a526035b3b0e66eddf559889dfae' \
      -X POST -d '{"pipelineId": "58813bb5bb0e020100e9525e", "message":"Manually triggered"}' https://app.wercker.com/api/v3/runs/
