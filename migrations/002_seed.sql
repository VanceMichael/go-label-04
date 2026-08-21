INSERT INTO tenants(id,name,active,created_at) VALUES('tenant-demo','Demo Cargo Tenant',true,now()) ON CONFLICT DO NOTHING;
