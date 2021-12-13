# Firopow

Firopow is a variation of ProgPow that uses modified Ethash DAG parameters,
modified ProgPow parameters, and a variation on the ProgPow initialization
and finalization functions.

## Modified Ethash Parameters

  - `DatasetInitBytes`: 2^30 + 2^29
  - `DatasetParents`: 512
  - `EpochLength`: 1300

## Modified ProgPow Parameters

  - `PeriodLength`: 1

## Modified Progpow Functions

These are a bit more complex, but should be clear in [firopow.go](./firopow.go).