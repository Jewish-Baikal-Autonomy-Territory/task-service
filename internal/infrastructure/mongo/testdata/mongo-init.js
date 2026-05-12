db = db.getSiblingDB("admin")

db.createUser({
    user: "backend-test",
    pwd: "backend-test",
    roles: [
        { role: "readWrite", db: "test" }
    ],
})

db = db.getSiblingDB("test")

db.createCollection("tasks", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: [
                "_id",
                "owner_id",
                "title",
                "description",
                "priority",
                "status",
                "created_at",
                "updated_at",
            ],
            properties: {
                _id: {
                    bsonType: "binData",
                    description: "Unique task id used for identification.",
                },
                owner_id: {
                    bsonType: "binData",
                    description: "Id of the task owner registered in the system.",
                },
                group_id: {
                    bsonType: "binData",
                    description: "Id of the group that task belong to.",
                },
                title: {
                    bsonType: "string",
                    minLength: 1,
                    description: "Summary representing task description.",
                },
                description: {
                    bsonType: "string",
                    minLength: 1,
                    maxLength: 10000,
                    description: "Description representing the contents of the task."
                },
                location: {
                    bsonType: "object",
                    required: [
                        "type",
                        "coordinates",
                    ],
                    properties: {
                        type: {
                            bsonType: "string",
                            enum: [
                                "Point",
                            ],
                        },
                        coordinates: {
                            bsonType: "array",
                            minItems: 2,
                            maxItems: 2,
                            items: [
                                {
                                    bsonType: "double",
                                    minimum: -180,
                                    maximum: 180,
                                },
                                {
                                    bsonType: "double",
                                    minimum: -90,
                                    maximum: 90,
                                },
                            ],
                        },
                    },
                    additionalProperties: false,
                },
                priority: {
                    bsonType: "string",
                    enum: [
                        "low",
                        "medium",
                        "high",
                    ],
                    description: "Priority represents the importance of task for the user.",
                },
                status: {
                    bsonType: "string",
                    enum: [
                        "pending",
                        "completed",
                    ],
                    description: "Status represents the current task state in the system."
                },
                is_favorite: {
                    bsonType: "bool",
                    description: "Boolean flag representing the pinned task status.",
                },
                created_at: {
                    bsonType: "date",
                    description: "The date when the task was created in the system.",
                },
                updated_at: {
                    bsonType: "date",
                    description: "The date when the task entity got updated in the system.",
                },
                completed_at: {
                    bsonType: "date",
                    description: "The date when the task is considered completed."
                },
                deadline: {
                    bsonType: "date",
                    description: "The date before the task should be completed.",
                },
                purge_at: {
                    bsonType: "date",
                    description: "The date when the task is going to be automatically deleted.",
                }
            },
            additionalProperties: false
        },
    },
    validationAction: "error",
    validationLevel: "strict",
})

db.tasks.createIndex({"location": "2dsphere"})
