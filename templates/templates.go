package templates

import (
	_ "embed"
)

type TemplatesSet struct {
	Short string
	Long  string
}

//go:embed en-template-short.md
var shortStandardTemplateEn string

// var ShortStandardTemplate = `# {{NUMBER}}. {{TITLE}}

// Date: {{DATE}}

// ## Status

// {{DATE}} proposed

// ## Context

// The issue motivating this decision, and any context that influences or constrains the decision.

// ## Decision

// The change that we're proposing or have agreed to implement.

// ## Consequences

// What becomes easier or more difficult to do and any risks introduced by the change that will need to be mitigated.
// `

//go:embed en-template-long.md
var longStandardTemplateEn string

// var LongStandardTemplate = `# {{NUMBER}}. {{TITLE}}

// Date: {{DATE}}

// Deciders: [list everyone involved in the decision] <!-- optional -->
// Technical Story: [description | ticket/issue URL] <!-- optional -->
// Pull Request: [PR URL] <!-- optional -->

// ## Status

// {{DATE}} proposed

// ## Context and Problem Statement

// [Describe the context and problem statement, e.g., in free form using two to three sentences. You may want to articulate the problem in form of a question.]

// ## Decision Drivers <!-- optional -->

// * [driver 1, e.g., a force, facing concern, …]
// * [driver 2, e.g., a force, facing concern, …]
// * … <!-- numbers of drivers can vary -->

// ## Considered Options

// * [option 1]
// * [option 2]
// * [option 3]
// * … <!-- numbers of options can vary -->

// ## Decision Outcome

// Chosen option: "[option 1]", because [justification. e.g., only option, which meets k.o. criterion decision driver | which resolves force force | … | comes out best (see below)].

// ### Positive Consequences <!-- optional -->

// * [e.g., improvement of quality attribute satisfaction, follow-up decisions required, …]
// * …

// ### Negative Consequences <!-- optional -->

// * [e.g., compromising quality attribute, follow-up decisions required, …]
// * …

// ## Pros and Cons of the Options <!-- optional -->

// ### [option 1]

// [example | description | pointer to more information | …] <!-- optional -->

// * Good, because [argument a]
// * Good, because [argument b]
// * Bad, because [argument c]
// * … <!-- numbers of pros and cons can vary -->

// ### [option 2]

// [example | description | pointer to more information | …] <!-- optional -->

// * Good, because [argument a]
// * Good, because [argument b]
// * Bad, because [argument c]
// * … <!-- numbers of pros and cons can vary -->

// ### [option 3]

// [example | description | pointer to more information | …] <!-- optional -->

// * Good, because [argument a]
// * Good, because [argument b]
// * Bad, because [argument c]
// * … <!-- numbers of pros and cons can vary -->

// ## Links <!-- optional -->

// * [Link type] [Link to ADR] <!-- example: Refined by [ADR-0005](0005-example.md) -->
// * … <!-- numbers of links can vary -->
// `

var TemplatesLibrary = make(map[string]TemplatesSet)

func init() {
	TemplatesLibrary["en"] = TemplatesSet{Short: shortStandardTemplateEn, Long: longStandardTemplateEn}
}
