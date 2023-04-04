#!/bin/bash
kubectl exec -i $(kubectl get pods | awk '{print $1}' | grep -e 'labs-sql') -- psql -U labs -d taskflows < labs-dump.sql
