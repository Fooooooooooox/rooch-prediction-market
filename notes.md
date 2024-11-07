
# Health Check
curl https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/health

# Health Check with Database
curl https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/healthDb

# Create a Market
curl -X POST https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/markets \
-H "Content-Type: application/json" \
-d '{"title": "Market Title", "description": "Market Description"}'


# Update a Market
curl -X PUT https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/markets/1 \
-H "Content-Type: application/json" \
-d '{"title": "Updated Market Title", "description": "Updated Market Description"}'

# Get all Markets
curl https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/markets

# Get Trades by Market ID
curl https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/trades/1

# Create a Trade
curl -X POST https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/trade \
-H "Content-Type: application/json" \
-d '{"address": "0x123", "market_id": 1, "side": "buy", "tick": "yes", "amount": 100}'

# Get Votes by Market ID
curl https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/votes/1

# Create a Vote
curl -X POST https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/vote \
-H "Content-Type: application/json" \
-d '{"address": "0x123", "market_id": 1, "tick": "yes", "sig": "signature", "amount": 50}'

# Add User Balance
curl -X POST https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/user-balance/add \
     -H "Content-Type: application/json" \
     -d '{
           "address": "user_address_1",
           "amount": 100
         }'

# Decrease User Balance
curl -X POST https://p01--rooch-prediction-market-backend--jqzlbjwpw4qd.code.run/api/v1/user-balance/decrease \
     -H "Content-Type: application/json" \
     -d '{
           "address": "user_address_1",
           "amount": 100
         }'
