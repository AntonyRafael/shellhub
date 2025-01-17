package mongo

import (
	"context"
	"sort"
	"testing"

	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/pkg/fixtures"
	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/cache"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFirewallRuleList(t *testing.T) {
	type Expected struct {
		rules []models.FirewallRule
		len   int
		err   error
	}

	cases := []struct {
		description string
		page        paginator.Query
		fixtures    []string
		expected    Expected
	}{
		{
			description: "succeeds when no firewall rules are found",
			page:        paginator.Query{Page: -1, PerPage: -1},
			fixtures:    []string{},
			expected: Expected{
				rules: []models.FirewallRule{},
				len:   0,
				err:   nil,
			},
		},
		{
			description: "succeeds when a firewall rule is found",
			page:        paginator.Query{Page: -1, PerPage: -1},
			fixtures:    []string{fixtures.FixtureFirewallRules},
			expected: Expected{
				rules: []models.FirewallRule{
					{
						ID:       "6504b7bd9b6c4a63a9ccc053",
						TenantID: "00000000-0000-4000-0000-000000000000",
						FirewallRuleFields: models.FirewallRuleFields{
							Priority: 1,
							Action:   "allow",
							Active:   true,
							SourceIP: ".*",
							Username: ".*",
							Filter: models.FirewallFilter{
								Hostname: "",
								Tags:     []string{"tag-1"},
							},
						},
					},
					{
						ID:       "e92f4a5d3e1a4f7b8b2b6e9a",
						TenantID: "00000000-0000-4000-0000-000000000000",
						FirewallRuleFields: models.FirewallRuleFields{
							Priority: 2,
							Action:   "allow",
							Active:   true,
							SourceIP: "192.168.1.10",
							Username: "john.doe",
							Filter: models.FirewallFilter{
								Hostname: "",
								Tags:     []string{"tag-1"},
							},
						},
					},
					{
						ID:       "78c96f0a2e5b4dca8d78f00c",
						TenantID: "00000000-0000-4000-0000-000000000000",
						FirewallRuleFields: models.FirewallRuleFields{
							Priority: 3,
							Action:   "allow",
							Active:   true,
							SourceIP: "10.0.0.0/24",
							Username: "admin",
							Filter: models.FirewallFilter{
								Hostname: "",
								Tags:     []string{},
							},
						},
					},
					{
						ID:       "3fd759a1ecb64ec5a07c8c0f",
						TenantID: "00000000-0000-4000-0000-000000000000",
						FirewallRuleFields: models.FirewallRuleFields{
							Priority: 4,
							Action:   "deny",
							Active:   true,
							SourceIP: "172.16.0.0/16",
							Username: ".*",
							Filter: models.FirewallFilter{
								Hostname: "",
								Tags:     []string{"tag-1"},
							},
						},
					},
				},
				len: 4,
				err: nil,
			},
		},
		{
			description: "succeeds when firewall rule list is not empty and paginator is different than -1",
			page:        paginator.Query{Page: 2, PerPage: 2},
			fixtures:    []string{fixtures.FixtureFirewallRules},
			expected: Expected{
				rules: []models.FirewallRule{
					{
						ID:       "78c96f0a2e5b4dca8d78f00c",
						TenantID: "00000000-0000-4000-0000-000000000000",
						FirewallRuleFields: models.FirewallRuleFields{
							Priority: 3,
							Action:   "allow",
							Active:   true,
							SourceIP: "10.0.0.0/24",
							Username: "admin",
							Filter: models.FirewallFilter{
								Hostname: "",
								Tags:     []string{},
							},
						},
					},
					{
						ID:       "3fd759a1ecb64ec5a07c8c0f",
						TenantID: "00000000-0000-4000-0000-000000000000",
						FirewallRuleFields: models.FirewallRuleFields{
							Priority: 4,
							Action:   "deny",
							Active:   true,
							SourceIP: "172.16.0.0/16",
							Username: ".*",
							Filter: models.FirewallFilter{
								Hostname: "",
								Tags:     []string{"tag-1"},
							},
						},
					},
				},
				len: 4,
				err: nil,
			},
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	// Due to the non-deterministic order of applying fixtures when dealing with multiple datasets,
	// we ensure that both the expected and result arrays are correctly sorted.
	sort := func(fr []models.FirewallRule) {
		sort.Slice(fr, func(i, j int) bool {
			return fr[i].ID < fr[j].ID
		})
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			rules, count, err := mongostore.FirewallRuleList(context.TODO(), tc.page)
			sort(tc.expected.rules)
			sort(rules)
			assert.Equal(t, tc.expected, Expected{rules: rules, len: count, err: err})
		})
	}
}

func TestFirewallRuleGet(t *testing.T) {
	type Expected struct {
		rule *models.FirewallRule
		err  error
	}
	cases := []struct {
		description string
		id          string
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when firewall rule is not found",
			id:          "6504b7bd9b6c4a63a9ccc021",
			fixtures:    []string{fixtures.FixtureFirewallRules},
			expected: Expected{
				rule: nil,
				err:  store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when firewall rule is found",
			id:          "6504b7bd9b6c4a63a9ccc053",
			fixtures:    []string{fixtures.FixtureFirewallRules},
			expected: Expected{
				rule: &models.FirewallRule{
					ID:       "6504b7bd9b6c4a63a9ccc053",
					TenantID: "00000000-0000-4000-0000-000000000000",
					FirewallRuleFields: models.FirewallRuleFields{
						Priority: 1,
						Action:   "allow",
						Active:   true,
						SourceIP: ".*",
						Username: ".*",
						Filter: models.FirewallFilter{
							Hostname: "",
							Tags:     []string{"tag-1"},
						},
					},
				},
				err: nil,
			},
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			rule, err := mongostore.FirewallRuleGet(context.TODO(), tc.id)
			assert.Equal(t, tc.expected, Expected{rule: rule, err: err})
		})
	}
}

func TestFirewallRuleUpdate(t *testing.T) {
	type Expected struct {
		rule *models.FirewallRule
		err  error
	}

	cases := []struct {
		description string
		id          string
		rule        models.FirewallRuleUpdate
		fixtures    []string
		expected    Expected
	}{
		{
			description: "fails when firewall rule is not found",
			id:          "6504b7bd9b6c4a63a9ccc000",
			rule: models.FirewallRuleUpdate{
				FirewallRuleFields: models.FirewallRuleFields{
					Priority: 1,
					Action:   "deny",
					Active:   true,
					SourceIP: ".*",
					Username: ".*",
					Filter: models.FirewallFilter{
						Hostname: "",
						Tags:     []string{"editedtag"},
					},
				},
			},
			fixtures: []string{fixtures.FixtureFirewallRules},
			expected: Expected{
				rule: nil,
				err:  store.ErrNoDocuments,
			},
		},
		{
			description: "succeeds when firewall rule is found",
			id:          "6504b7bd9b6c4a63a9ccc053",
			rule: models.FirewallRuleUpdate{
				FirewallRuleFields: models.FirewallRuleFields{
					Priority: 1,
					Action:   "deny",
					Active:   true,
					SourceIP: ".*",
					Username: ".*",
					Filter: models.FirewallFilter{
						Hostname: "",
						Tags:     []string{"editedtag"},
					},
				},
			},
			fixtures: []string{fixtures.FixtureFirewallRules},
			expected: Expected{
				rule: &models.FirewallRule{
					ID:       "6504b7bd9b6c4a63a9ccc053",
					TenantID: "00000000-0000-4000-0000-000000000000",
					FirewallRuleFields: models.FirewallRuleFields{
						Priority: 1,
						Action:   "deny",
						Active:   true,
						SourceIP: ".*",
						Username: ".*",
						Filter: models.FirewallFilter{
							Hostname: "",
							Tags:     []string{"editedtag"},
						},
					},
				},
				err: nil,
			},
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			rule, err := mongostore.FirewallRuleUpdate(context.TODO(), tc.id, tc.rule)
			assert.Equal(t, tc.expected, Expected{rule: rule, err: err})
		})
	}
}

func TestFirewallRuleDelete(t *testing.T) {
	cases := []struct {
		description string
		id          string
		fixtures    []string
		expected    error
	}{
		{
			description: "fails when rule is not found",
			id:          "6504ac006bf3dbca079f76b1",
			fixtures:    []string{},
			expected:    store.ErrNoDocuments,
		},
		{
			description: "succeeds when rule is found",
			id:          "6504b7bd9b6c4a63a9ccc053",
			fixtures:    []string{fixtures.FixtureFirewallRules},
			expected:    nil,
		},
	}

	db := dbtest.DBServer{}
	defer db.Stop()

	mongostore := NewStore(db.Client().Database("test"), cache.NewNullCache())
	fixtures.Init(db.Host, "test")

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.NoError(t, fixtures.Apply(tc.fixtures...))
			defer fixtures.Teardown() // nolint: errcheck

			err := mongostore.FirewallRuleDelete(context.TODO(), tc.id)
			assert.Equal(t, tc.expected, err)
		})
	}
}
