package commands

// MySkills is the list of skills you want installed on your machine
// The "Install All" action will install these skills
var MySkills = []struct {
	Repo  string // Repository in owner/repo format
	Skill string // Skill name
}{
	{"anthropics/skills", "frontend-design"},
	{"expo/skills", "upgrading-expo"},
	{"giuseppe-trisciuoglio/developer-kit", "shadcn-ui"},
	{"sickn33/antigravity-awesome-skills", "last30days"},
	{"vercel-labs/agent-skills", "vercel-react-best-practices"},
	{"vercel-labs/agent-skills", "vercel-react-native-skills"},
}

// SkillRepo represents a repository containing skills
type SkillRepo struct {
	Name        string // Repository name (owner/repo format)
	Description string // Brief description
}

// SkillRepos is the list of recommended skill repositories
// Skills are fetched dynamically when a repo is expanded in the TUI
var SkillRepos = []SkillRepo{
	{"anthropics/skills", "Official Anthropic skills for Claude"},
	{"better-auth/skills", "Authentication best practices"},
	{"code-with-beto/skills", "Beto's development skills"},
	{"coreyhaines31/marketingskills", "Marketing and SEO skills"},
	{"expo/skills", "Expo and React Native mobile development"},
	{"firecrawl/cli", "Web content extraction for AI agents"},
	{"giuseppe-trisciuoglio/developer-kit", "Developer toolkit including shadcn-ui"},
	{"obra/superpowers", "Development workflow and productivity skills"},
	{"remotion-dev/skills", "Remotion video creation skills"},
	{"resend/email-best-practices", "Email development best practices"},
	{"supabase/agent-skills", "Supabase database and backend skills"},
	{"vercel-labs/agent-skills", "Vercel React and web development skills"},
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
