package middleware

import (
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/rooch-prediction-market/backend/config"
	"github.com/rooch-prediction-market/backend/handlers"
	"gopkg.in/macaron.v1"
)

func ClerkMiddleware(c *macaron.Context) {
	clerkMiddleware := clerkhttp.WithHeaderAuthorization()
	req := c.Req.Request
	respWriter := c.Resp

	clerkMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			c.JSON(http.StatusUnauthorized, "Unauthorized: Clerk session not found")
			return
		}

		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error fetching user: %v", err))
			return
		}
		if usr == nil {
			c.JSON(http.StatusNotFound, "User does not exist")
			return
		}

		fmt.Println("this is UserId: ", usr.ID)
		c.Data["userClerkID"] = usr.ID

		c.Next()
	})).ServeHTTP(respWriter, req)
}

func ClerkKycReviewerMiddleware(c *macaron.Context) {
	clerkMiddleware := clerkhttp.WithHeaderAuthorization()
	req := c.Req.Request
	respWriter := c.Resp

	clerkMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fmt.Println("===this is clerk kyc reviewer middleware")
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			fmt.Println("===failed to get clerk claims", ok)
			c.JSON(http.StatusUnauthorized, "Unauthorized: Clerk session not found")
			return
		}

		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			fmt.Println("===failed to get user", err)
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error fetching user: %v", err))
			return
		}
		if usr == nil {
			fmt.Println("===user does not exist")
			c.JSON(http.StatusNotFound, "User does not exist")
			return
		}

		memberships, err := user.ListOrganizationMemberships(ctx, usr.ID, &user.ListOrganizationMembershipsParams{
			ListParams: clerk.ListParams{
				Limit:  clerk.Int64(10),
				Offset: clerk.Int64(0),
			},
		})

		len := int(memberships.TotalCount)
		if len >= 1 {
			for i := 0; i < len; i++ {
				fmt.Println("===membership role: ", memberships.OrganizationMemberships[i].Role)
				fmt.Println("===membership org name: ", memberships.OrganizationMemberships[i].Organization.Name)
				if memberships.OrganizationMemberships[i].Role == "org:admin" && memberships.OrganizationMemberships[i].Organization.Name == "kyc reviewer" {
					c.Data["userClerkID"] = usr.ID
					c.Next()
					return
				}
			}
			fmt.Println("===user is not a kyc reviewer with admin role")
			c.JSON(http.StatusForbidden, "User is not a kyc reviewer with admin role")
		} else {
			fmt.Println("===user is not a kyc reviewer")
			c.JSON(http.StatusForbidden, "User is not a kyc reviewer")
		}

		if err != nil {
			fmt.Println("===failed to fetch user memberships", err)
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error fetching user: %v", err))
			return
		}

		fmt.Println("this is UserId: ", usr.ID)
		c.Data["userClerkID"] = usr.ID

		c.Next()
	})).ServeHTTP(respWriter, req)
}

func JwtMiddleware(c *macaron.Context, config *config.Config) {
	authHeader := c.Req.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString := authHeader
	claims := &handlers.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// return handlers.JwtKey, nil
		return []byte(config.JwtKey), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, fmt.Sprintf("jwt verify failed, error: %s ", err))
		return
	}

	fmt.Println("this is UserId: ", claims.UserId)
	fmt.Println("this is UserTwitterId: ", claims.UserTwitterId)

	c.Data["userid"] = claims.UserId
	c.Data["userTwitterId"] = claims.UserTwitterId
	c.Next()
}
