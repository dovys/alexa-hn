package stub

import "github.com/dovys/alexa-hn/hn"

func NewStubHNClient() hn.Client {
	c := &stubHnClient{}

	c.topStories = []int{14779881, 14780709, 14778977, 14778685, 14780159, 14779509, 14778335}
	c.stories = map[int]*hn.Story{
		14779881: &hn.Story{
			ID:    14779881,
			Score: 571,
			Title: "Apache Foundation bans projects from using React's “BSD+Patent” Code",
		},
		14780709: &hn.Story{
			ID:    14780709,
			Score: 35,
			Title: "The Deal on the Table (1994)",
		},
		14778977: &hn.Story{
			ID:    14778977,
			Score: 177,
			Title: "A 32-year-old state senator is trying to get patent trolls out of Mass (techcrunch.com)",
		},
		14778685: &hn.Story{
			ID:    14778685,
			Score: 279,
			Title: "Monolith First (2015)",
		},
		14780159: &hn.Story{
			ID:    14780159,
			Score: 71,
			Title: "A deep dive into Multicore OCaml garbage collector",
		},
		14779509: &hn.Story{
			ID:    14779509,
			Score: 101,
			Title: "Why MAC address randomization is not enough [pdf]",
		},
		14778335: &hn.Story{
			ID:    14778335,
			Score: 207,
			Title: "Tokyo street fashion and culture: 1980 – 2017 ",
		},
	}

	return c
}

type stubHnClient struct {
	topStories []int
	stories    map[int]*hn.Story
}

func (c *stubHnClient) TopStories() ([]int, error) {
	return c.topStories, nil
}

func (c *stubHnClient) Story(id int) (*hn.Story, error) {
	return c.stories[id], nil
}
