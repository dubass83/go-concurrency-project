-- Drop foregin key constraints first
ALTER TABLE public.user_plans
DROP CONSTRAINT IF EXISTS user_plans_plan_id_fkey;

ALTER TABLE public.user_plans
DROP CONSTRAINT IF EXISTS user_plans_user_id_fkey;

-- Drop primary key constraints
ALTER TABLE public.users
DROP CONSTRAINT IF EXISTS user_pkey;

ALTER TABLE public.user_plans
DROP CONSTRAINT IF EXISTS user_plan_pkey;

ALTER TABLE public.plans
DROP CONSTRAINT IF EXISTS plans_pkey;

-- Drop tables
DROP TABLE IF EXISTS public.user_plans;

DROP TABLE IF EXISTS public.users;

DROP TABLE IF EXISTS public.plans;

-- DROP sequences
DROP SEQUENCE IF EXISTS public.user_plans_id_seq;

DROP SEQUENCE IF EXISTS public.user_id_seq;

DROP SEQUENCE IF EXISTS public.plans_id_seq;
