#!/bin/bash

# RabbitMQ Initialization Script
# Creates exchanges, queues, and bindings for event-driven architecture

echo "Initializing RabbitMQ..."

# Wait for RabbitMQ to be ready
sleep 5

# Create exchanges
rabbitmqctl declare_exchange name=tunduck.auth kind=topic durable=true
rabbitmqctl declare_exchange name=tunduck.company kind=topic durable=true
rabbitmqctl declare_exchange name=tunduck.document kind=topic durable=true

echo "✅ Exchanges created:"
echo "   - tunduck.auth (topic exchange)"
echo "   - tunduck.company (topic exchange)"
echo "   - tunduck.document (topic exchange)"

# Create queues for Auth Service events
rabbitmqctl declare_queue name=auth.events.user_registered durable=true
rabbitmqctl declare_queue name=auth.events.user_logged_in durable=true
rabbitmqctl declare_queue name=auth.events.user_logged_out durable=true
rabbitmqctl declare_queue name=auth.events.user_registered.dlx durable=true

echo "✅ Auth Service queues created"

# Create queues for Company Service events
rabbitmqctl declare_queue name=company.events.organization_created durable=true
rabbitmqctl declare_queue name=company.events.organization_updated durable=true
rabbitmqctl declare_queue name=company.events.organization_deleted durable=true
rabbitmqctl declare_queue name=company.events.employee_added durable=true
rabbitmqctl declare_queue name=company.events.employee_removed durable=true
rabbitmqctl declare_queue name=company.events.organization_created.dlx durable=true

echo "✅ Company Service queues created"

# Create queues for Document Service events
rabbitmqctl declare_queue name=document.events.document_created durable=true
rabbitmqctl declare_queue name=document.events.document_sent durable=true
rabbitmqctl declare_queue name=document.events.document_approved durable=true
rabbitmqctl declare_queue name=document.events.document_rejected durable=true
rabbitmqctl declare_queue name=document.events.document_archived durable=true
rabbitmqctl declare_queue name=document.events.document_created.dlx durable=true

echo "✅ Document Service queues created"

# Bind exchanges to queues (routing)

# Auth Service bindings
rabbitmqctl bind_queue exchange=tunduck.auth queue=auth.events.user_registered routing_key="user.registered"
rabbitmqctl bind_queue exchange=tunduck.auth queue=auth.events.user_logged_in routing_key="user.logged_in"
rabbitmqctl bind_queue exchange=tunduck.auth queue=auth.events.user_logged_out routing_key="user.logged_out"

echo "✅ Auth Service bindings created"

# Company Service bindings (listening to auth events)
rabbitmqctl bind_queue exchange=tunduck.auth queue=company.events.organization_created routing_key="user.registered"

# Company Service publishes its own events
rabbitmqctl bind_queue exchange=tunduck.company queue=company.events.organization_created routing_key="organization.created"
rabbitmqctl bind_queue exchange=tunduck.company queue=company.events.organization_updated routing_key="organization.updated"
rabbitmqctl bind_queue exchange=tunduck.company queue=company.events.organization_deleted routing_key="organization.deleted"
rabbitmqctl bind_queue exchange=tunduck.company queue=company.events.employee_added routing_key="employee.added"
rabbitmqctl bind_queue exchange=tunduck.company queue=company.events.employee_removed routing_key="employee.removed"

echo "✅ Company Service bindings created"

# Document Service bindings (listening to auth and company events)
rabbitmqctl bind_queue exchange=tunduck.auth queue=document.events.document_created routing_key="user.registered"
rabbitmqctl bind_queue exchange=tunduck.company queue=document.events.document_created routing_key="organization.created"

# Document Service publishes its own events
rabbitmqctl bind_queue exchange=tunduck.document queue=document.events.document_created routing_key="document.created"
rabbitmqctl bind_queue exchange=tunduck.document queue=document.events.document_sent routing_key="document.sent"
rabbitmqctl bind_queue exchange=tunduck.document queue=document.events.document_approved routing_key="document.approved"
rabbitmqctl bind_queue exchange=tunduck.document queue=document.events.document_rejected routing_key="document.rejected"
rabbitmqctl bind_queue exchange=tunduck.document queue=document.events.document_archived routing_key="document.archived"

echo "✅ Document Service bindings created"

# Setup Dead Letter Exchange for error handling
rabbitmqctl declare_exchange name=tunduck.dlx kind=topic durable=true

# Bind dead letter queues
rabbitmqctl bind_queue exchange=tunduck.dlx queue=auth.events.user_registered.dlx routing_key="*.dlx"
rabbitmqctl bind_queue exchange=tunduck.dlx queue=company.events.organization_created.dlx routing_key="*.dlx"
rabbitmqctl bind_queue exchange=tunduck.dlx queue=document.events.document_created.dlx routing_key="*.dlx"

echo "✅ Dead Letter Exchange setup completed"

# Set queue arguments for dead letter routing
rabbitmqctl set_queue_attributes queue=auth.events.user_registered \
  deadletterexchange=tunduck.dlx \
  deadletterroutingkey="auth.dlx" \
  "x-message-ttl"=3600000 \
  "x-max-retries"=3

echo "✅ Queue attributes configured"

echo "🎉 RabbitMQ initialization completed successfully!"
echo ""
echo "Created exchanges:"
echo "  - tunduck.auth"
echo "  - tunduck.company"
echo "  - tunduck.document"
echo "  - tunduck.dlx (dead letter)"
echo ""
echo "Management UI: http://localhost:15672"
echo "Default credentials: tunduck_user / tunduck_password_dev"
