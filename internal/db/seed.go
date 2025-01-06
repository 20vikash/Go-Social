package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"social/social/internal/store"
)

var usernames = []string{
	"Alice", "Bob", "Charlie", "Dave", "Eve",
	"Frank", "Grace", "Hannah", "Ivy", "Jack",
	"Karen", "Liam", "Mona", "Nate", "Olivia",
	"Paul", "Quincy", "Rachel", "Sam", "Tina",
	"Uma", "Victor", "Wendy", "Xander", "Yara",
	"Zane", "Anna", "Ben", "Clara", "Derek",
	"Emma", "Felix", "Georgia", "Harry", "Isla",
	"James", "Katie", "Lucas", "Megan", "Noah",
	"Oscar", "Penny", "Quinn", "Ryan", "Sophie",
	"Tom", "Ursula", "Violet", "Will", "Zoe",
}

var titles = []string{
	"5 Tips for Better Time Management",
	"How to Stay Productive While Working Remotely",
	"The Beginner's Guide to Investing",
	"10 Healthy Snacks for Busy People",
	"Top 5 Benefits of Daily Exercise",
	"The Art of Minimalist Living",
	"How to Learn a New Skill in 30 Days",
	"The Future of Artificial Intelligence",
	"Traveling on a Budget: Hacks You Need to Know",
	"Why Reading Books is Good for Your Brain",
	"The Power of Gratitude in Everyday Life",
	"How to Start a Successful Blog in 2025",
	"Understanding Blockchain in Simple Terms",
	"Tips for Maintaining a Healthy Work-Life Balance",
	"The Rise of Electric Vehicles: What You Should Know",
	"How to Create a Weekly Meal Plan That Works",
	"The Importance of Mental Health Awareness",
	"10 Simple Ways to Reduce Your Carbon Footprint",
	"How to Master Public Speaking",
	"The Benefits of Journaling for Personal Growth",
}

var contents = []string{
	"Time management is a critical skill that allows you to balance work, personal life, and leisure. Start by creating a priority list and stick to it daily.",
	"Working remotely can be productive if you set boundaries and create a dedicated workspace. Don't forget to take regular breaks to avoid burnout.",
	"Investing might seem complex, but starting with index funds or ETFs can make it simple and effective. Research is key to understanding the market.",
	"When you're busy, healthy snacks like nuts, fruits, or yogurt can keep your energy levels up. Always have them handy for a quick bite.",
	"Exercising daily boosts not just your physical health but also your mental well-being. A simple 30-minute walk can make a big difference.",
	"Minimalism is more than owning fewer things; it's about creating space for what truly matters in your life. Start by decluttering one area of your home.",
	"Learning a new skill doesn't have to be daunting. Break it into manageable tasks and practice consistently to see improvement over time.",
	"Artificial Intelligence is transforming industries like healthcare, finance, and education. The potential for AI is vast, but ethical considerations remain crucial.",
	"Traveling doesn't have to be expensive. Use comparison websites, travel during off-peak seasons, and consider staying in hostels or Airbnbs to save money.",
	"Reading stimulates your imagination and strengthens your critical thinking skills. Make it a habit to read at least one book per month.",
	"Gratitude shifts your focus from what's missing in your life to what you already have. A daily gratitude journal can improve your perspective and mood.",
	"Starting a blog is easier than ever. Choose a niche, create a content calendar, and use social media to promote your posts effectively.",
	"Blockchain technology ensures secure and transparent transactions. While often associated with cryptocurrencies, its applications extend to supply chain and healthcare.",
	"Achieving work-life balance requires setting boundaries, delegating tasks, and prioritizing self-care. Remember, rest is productive too.",
	"Electric vehicles are becoming more accessible and eco-friendly. With advancements in charging infrastructure, EVs are set to dominate the future of transportation.",
	"Meal planning saves time and reduces food waste. Start by planning your meals for the week, making a shopping list, and sticking to it.",
	"Mental health is as important as physical health. Seek support, practice self-care, and educate yourself about mental health challenges to help others too.",
	"Reducing your carbon footprint can be as simple as using reusable bags, cutting down on single-use plastics, and carpooling whenever possible.",
	"Public speaking is a skill that can be mastered with practice. Start small, rehearse your speeches, and focus on engaging your audience.",
	"Journaling helps you process your thoughts, set goals, and track personal growth. Even five minutes a day can have a profound impact on your mindset.",
}

var tags = []string{
	"Productivity", "Health", "Technology", "Travel", "Finance",
	"Self-Improvement", "Minimalism", "AI", "Education", "Lifestyle",
	"Sustainability", "Fitness", "Mental Health", "Blogging", "Blockchain",
	"Cooking", "Time Management", "Public Speaking", "Personal Growth", "Reading",
}

var commentsD = []string{
	"Great insights! This was really helpful, thanks for sharing.",
	"I never thought about it this way. Totally agree with your point.",
	"This is exactly what I needed to read today. Thank you!",
	"I have been struggling with this issue for a while. Your tips make so much sense.",
	"Interesting perspective! I'd love to hear more about this topic.",
	"Fantastic article! Could you provide some examples to go along with this?",
	"I've been using some of these strategies, and they really work.",
	"This was super informative. Keep up the great work!",
	"I've bookmarked this to revisit later. Thanks for putting this together!",
	"Do you have any recommendations for beginners in this field?",
	"This resonated with me deeply. Thanks for shedding light on this issue.",
	"I appreciate the clear and concise explanations here.",
	"I shared this with my friends, and they found it helpful too.",
	"Looking forward to more content like this from you!",
	"Could you expand on this point? It's really intriguing.",
	"This was an eye-opener. Thanks for the detailed write-up.",
	"Your writing style is engaging and easy to follow. Loved it!",
	"I've learned something new today. Thank you for this!",
	"Can you recommend any resources or books on this topic?",
	"This sparked a lot of ideas for me. Appreciate your efforts!",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", post)
			return
		}
	}

	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", comment)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]

		comments[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: commentsD[rand.Intn(len(commentsD))],
		}
	}

	return comments
}
