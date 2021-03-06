package templates

// OpeningHTML is the opening portion html tag
const OpeningHTML = `
<!DOCTYPE html5>
<html>
`

// ClosingHTML is the closing portion html tag
const ClosingHTML = `
</html>
`

// Head is the head portion of the html template
const Head = `
<head>
	<title>GoJournalHTTP</title>
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>
</head>
`

// OpeningBody is the opening portion of the body tag
const OpeningBody = `
<body>
	<div class="container">
		<div class="row">
			<div class="col">
				<h1 style="text-align:center;">Welcome to GoJournal!</h1>
			</div>
		</div>
		<div class="row">
			<div class="col"></div>
			<div class="col">
			</div>
			<div class="col"></div>
		</div>
		<div class="row">
			<div class="col"></div>
			<div class="col">
				<form action="/" method="POST">
					<input type="text" name="date" placeholder="YYYY-MM-DD"/>
					<textarea rows="8" cols="50" type="text" name="entry" placeholder="Type your journal entry"></textarea>
					<button type="submit">Add Entry</button>
				</form>
`

// ClosingBody is the closing portion of the body tag
const ClosingBody = `
			</div>
			<div class="col"></div>
		</div>
	</div>
</body>
`
