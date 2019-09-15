from SPARQLWrapper import SPARQLWrapper, JSON
import pywikibot
import pandas as pd
import os
site = pywikibot.Site("wikidata", "wikidata")
repo = site.data_repository()

f_query_name="fetchstore.txt"
f_property_name="fetchprops.txt"

if not os.path.exists(f_query_name):
	sparql = SPARQLWrapper("https://query.wikidata.org/sparql")
	sparql.setQuery("""
	SELECT ?item ?itemLabel 
	WHERE 
	{
	  ?item wdt:P31 wd:Q16521.
	  ?item wdt:P171* wd:Q5113.
	  ?item wdt:P105 wd:Q7432.
	  SERVICE wikibase:label { bd:serviceParam wikibase:language "[AUTO_LANGUAGE],ru". }
	}""")
	sparql.setReturnFormat(JSON)
	results = sparql.query().convert()
	results_df = pd.io.json.json_normalize(results['results']['bindings'])
	keep=results_df[['item.value', 'itemLabel.value']].values.tolist()
	p='\n'.join(['\t'.join(k) for k in keep])
	with open(f_query_name,"w",encoding="utf-8") as f:
		f.write(p)

with open(f_query_name,encoding="utf-8") as f:
	fetched=f.read()

def qq(qid):
    item = pywikibot.ItemPage(repo, qid)
    item_dict = item.get()
    clm_dict = item_dict["claims"]
    return list(clm_dict.keys())

v=fetched.split('\n')
d=[v0.split('\t') for v0 in v]
Qs=[i[0].split('/')[-1] for i in d]
Ps=[[q,qq(q)] for q in Qs]
p='\n'.join(['\t\t'.join((k[0],'\t'.join(k[1]))) for k in Ps])
with open(f_property_name,"w",encoding="utf-8") as f:
	f.write(p)

