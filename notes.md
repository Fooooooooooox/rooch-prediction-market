# Health Check
curl http://localhost:50051/api/v1/health

# Health Check with Database
curl http://localhost:50051/api/v1/healthDb

# Add User Balance
curl -X POST http://localhost:50051/api/v1/user-balance/add \
     -H "Content-Type: application/json" \
     -d '{"address": "user_address_bcbc", "amount": 100}'

# Decrease User Balance
curl -X POST http://localhost:50051/api/v1/user-balance/decrease \
     -H "Content-Type: application/json" \
     -d '{"address": "user_address_bcbc", "amount": 100}'

# Get User Balance
curl http://localhost:50051/api/v1/user-balance/user_address_bcbc

# Get User Market Balance
curl http://localhost:50051/api/v1/user-market-balance/user_address_bcbc/2

# Create a Market
curl -X POST http://localhost:50051/api/v1/markets \
     -H "Content-Type: application/json" \
     -d '{"title": "Market Title", "description": "Market Description"}'

# Get all Markets
curl http://localhost:50051/api/v1/markets

# Get a Market by ID
curl http://localhost:50051/api/v1/markets/1

# Update a Market
curl -X PUT http://localhost:50051/api/v1/markets/update \
     -H "Content-Type: application/json" \
     -d '{"market_id": 1, "title": "Updated Market Title", "description": "Updated Market Description"}'

# Settle a Market
curl -X POST http://localhost:50051/api/v1/markets/settle \
     -H "Content-Type: application/json" \
     -d '{"market_id": 1, "result": true}'

# Get Trades by Market ID
curl http://localhost:50051/api/v1/trades/1

# Create a Trade
curl -X POST http://localhost:50051/api/v1/trade \
     -H "Content-Type: application/json" \
     -d '{"address": "user_address_bcbc", "market_id": 1, "side": "buy", "tick": "yes", "amount": 100}'

# Get Votes by Market ID
curl http://localhost:50051/api/v1/votes/1

# Create a Vote
curl -X POST http://localhost:50051/api/v1/vote \
     -H "Content-Type: application/json" \
     -d '{"address": "user_address_bcbc", "market_id": 1, "tick": "yes", "amount": 50}'

# Get Claimable Amount
curl http://localhost:50051/api/v1/claim/user_address_bcbc/1

# Claim Reward
curl -X POST http://localhost:50051/api/v1/claim \
     -H "Content-Type: application/json" \
     -d '{"address": "user_address_bcbc", "market_id": 1}'