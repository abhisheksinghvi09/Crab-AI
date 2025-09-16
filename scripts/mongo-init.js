// MongoDB initialization script for development environment
db = db.getSiblingDB('crab-ai');

// Create collections
db.createCollection('users');
db.createCollection('chats');
db.createCollection('llm_messages');
db.createCollection('waitlist');

// Create indexes for better performance
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "created_at": 1 });
db.chats.createIndex({ "user_id": 1 });
db.chats.createIndex({ "created_at": 1 });
db.llm_messages.createIndex({ "chat_id": 1 });
db.llm_messages.createIndex({ "created_at": 1 });
db.waitlist.createIndex({ "email": 1 }, { unique: true });

print('Database initialization completed successfully!');