INSERT INTO
    user_ (id, email, pass_hash, pass_salt, name_)
VALUES
    -- Admins
    ('b6c6cd06-e367-4a85-bf73-176da944a8c0', 'admin@example.com', '$2a$10$N7ocLLFa0j3AHXxHTGdI4.FeDCJaRLI4K3SiOMofRnup076JR7Z5y', 'admin123', 'Admin'), 
    
    -- President users
    ('0943d66c-de3e-4650-b60b-c182e23d6e9f', 'president@example.com', '$2a$10$pjUGnrVfAGak2PcDT96nWOfNIb/qLaH5.cLMJTmJWxSZH.0Ogi04u', 'president123', 'President'),

    -- Regular users
    ('f4106262-7e8f-4b31-9a8d-e02d5545c9c5', 'user@example.com', '$2a$10$QzKOEgNM.PKbQ2Q5cqVVO.t62p3zeHPcfpo3G1hgRAuvO8pk0YB.O', 'user1234', 'User')
ON CONFLICT DO NOTHING;

INSERT INTO
    institution (id, name_)
VALUES 
    ('775f19a3-a2be-4561-88e0-d82f98ee2e6e', 'Institution 1One'),
    ('20d79ab4-8782-43c1-b992-6227ebe3723e', 'Institution 2Two')
ON CONFLICT DO NOTHING;

INSERT INTO
    member (id, category, institution_id, user_id)
VALUES
    ('ba66f545-ae26-4fcd-b48d-139c678f9570', 'PRESIDENT', '775f19a3-a2be-4561-88e0-d82f98ee2e6e', '0943d66c-de3e-4650-b60b-c182e23d6e9f'),
    ('f4301d55-6354-43e1-bd3f-bddace5fbb7c', 'USER', '775f19a3-a2be-4561-88e0-d82f98ee2e6e', 'f4106262-7e8f-4b31-9a8d-e02d5545c9c5')
ON CONFLICT DO NOTHING;