# Kawpow

Kawpow is a variation of ProgPow that uses modified Ethash DAG parameters,
modified ProgPow parameters, and a variation on the ProgPow initialization
and finalization functions.

## Modified Ethash Parameters

  - `DatasetParents`: 512
  - `EpochLength`: 7500

## Modified ProgPow Parameters

  - `PeriodLength`: 3

## Modified Progpow Functions

These are a bit more complex, but should be clear in [kawpow.go](./kawpow.go).