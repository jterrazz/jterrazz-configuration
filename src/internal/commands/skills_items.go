package commands

// MySkills is the list of skills you want installed on your machine
// The "Install All" action will install these skills
var MySkills = []struct {
	Repo  string // Repository in owner/repo format
	Skill string // Skill name
}{
	{"anthropics/skills", "frontend-design"},
}

// SkillRepo represents a repository containing skills
type SkillRepo struct {
	Name        string   // Repository name (owner/repo format)
	Description string   // Brief description
	Skills      []string // List of skill names in this repo
}

// SkillRepos is the list of recommended skill repositories
var SkillRepos = []SkillRepo{
	{
		Name:        "anthropics/skills",
		Description: "Official Anthropic skills for Claude",
		Skills: []string{
			"frontend-design",
			"skill-creator",
			"pdf",
			"pptx",
			"xlsx",
			"docx",
			"canvas-design",
			"mcp-builder",
			"web-artifacts-builder",
			"brand-guidelines",
			"internal-comms",
			"theme-factory",
			"doc-coauthoring",
			"slack-gif-creator",
			"template-skill",
		},
	},
	{
		Name:        "vercel-labs/agent-skills",
		Description: "Vercel React and web development skills",
		Skills: []string{
			"vercel-react-best-practices",
			"vercel-composition-patterns",
			"vercel-react-native-skills",
			"web-design-guidelines",
		},
	},
	{
		Name:        "expo/skills",
		Description: "Expo and React Native mobile development",
		Skills: []string{
			"building-native-ui",
			"upgrading-expo",
			"native-data-fetching",
			"expo-dev-client",
			"expo-deployment",
			"expo-tailwind-setup",
			"expo-api-routes",
			"expo-cicd-workflows",
			"use-dom",
		},
	},
	{
		Name:        "obra/superpowers",
		Description: "Development workflow and productivity skills",
		Skills: []string{
			"brainstorming",
			"systematic-debugging",
			"test-driven-development",
			"writing-plans",
			"executing-plans",
			"subagent-driven-development",
			"verification-before-completion",
			"requesting-code-review",
			"using-superpowers",
			"writing-skills",
			"dispatching-parallel-agents",
			"using-git-worktrees",
			"receiving-code-review",
			"finishing-a-development-branch",
		},
	},
	{
		Name:        "coreyhaines31/marketingskills",
		Description: "Marketing and SEO skills",
		Skills: []string{
			"seo-audit",
			"copywriting",
			"marketing-psychology",
			"programmatic-seo",
			"pricing-strategy",
			"social-content",
			"copy-editing",
			"launch-strategy",
			"page-cro",
			"analytics-tracking",
			"onboarding-cro",
			"schema-markup",
			"competitor-alternatives",
			"paid-ads",
			"email-sequence",
		},
	},
	{
		Name:        "supabase/agent-skills",
		Description: "Supabase database and backend skills",
		Skills: []string{
			"supabase-postgres-best-practices",
		},
	},
	{
		Name:        "remotion-dev/skills",
		Description: "Remotion video creation skills",
		Skills: []string{
			"remotion-best-practices",
		},
	},
	{
		Name:        "better-auth/skills",
		Description: "Authentication best practices",
		Skills: []string{
			"better-auth-best-practices",
			"create-auth-skill",
		},
	},
}

// GetAllSkillRepos returns all skill repositories
func GetAllSkillRepos() []SkillRepo {
	return SkillRepos
}

// GetSkillRepoByName returns a skill repo by name
func GetSkillRepoByName(name string) *SkillRepo {
	for i := range SkillRepos {
		if SkillRepos[i].Name == name {
			return &SkillRepos[i]
		}
	}
	return nil
}
