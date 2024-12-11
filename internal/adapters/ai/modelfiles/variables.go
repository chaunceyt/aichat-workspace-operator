/*
Copyright 2024 AIChatWorkspace Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package modelfiles

import (
	_ "embed"
)

//go:embed files/agility_story/system.md
var agility_story string

//go:embed files/ai/system.md
var ai string

//go:embed files/analyze_answers/system.md
var analyze_answers string

//go:embed files/analyze_candidates/system.md
var analyze_candidates string

//go:embed files/analyze_cfp_submission/system.md
var analyze_cfp_submission string

//go:embed files/analyze_claims/system.md
var analyze_claims string

//go:embed files/analyze_comments/system.md
var analyze_comments string

//go:embed files/analyze_debate/system.md
var analyze_debate string

//go:embed files/analyze_email_headers/system.md
var analyze_email_headers string

//go:embed files/analyze_incident/system.md
var analyze_incident string

//go:embed files/analyze_interviewer_techniques/system.md
var analyze_interviewer_techniques string

//go:embed files/analyze_logs/system.md
var analyze_logs string

//go:embed files/analyze_malware/system.md
var analyze_malware string

//go:embed files/analyze_military_strategy/system.md
var analyze_military_strategy string

//go:embed files/analyze_mistakes/system.md
var analyze_mistakes string

//go:embed files/analyze_paper/system.md
var analyze_paper string

//go:embed files/analyze_patent/system.md
var analyze_patent string

//go:embed files/analyze_personality/system.md
var analyze_personality string

//go:embed files/analyze_presentation/system.md
var analyze_presentation string

//go:embed files/analyze_product_feedback/system.md
var analyze_product_feedback string

//go:embed files/analyze_proposition/system.md
var analyze_proposition string

//go:embed files/analyze_prose/system.md
var analyze_prose string

//go:embed files/analyze_prose_json/system.md
var analyze_prose_json string

//go:embed files/analyze_prose_pinker/system.md
var analyze_prose_pinker string

//go:embed files/analyze_risk/system.md
var analyze_risk string

//go:embed files/analyze_sales_call/system.md
var analyze_sales_call string

//go:embed files/analyze_spiritual_text/system.md
var analyze_spiritual_text string

//go:embed files/analyze_tech_impact/system.md
var analyze_tech_impact string

//go:embed files/analyze_threat_report/system.md
var analyze_threat_report string

//go:embed files/analyze_threat_report_trends/system.md
var analyze_threat_report_trends string

//go:embed files/answer_interview_question/system.md
var answer_interview_question string

//go:embed files/ask_secure_by_design_questions/system.md
var ask_secure_by_design_questions string

//go:embed files/ask_uncle_duke/system.md
var ask_uncle_duke string

//go:embed files/capture_thinkers_work/system.md
var capture_thinkers_work string

//go:embed files/check_agreement/system.md
var check_agreement string

//go:embed files/clean_text/system.md
var clean_text string

//go:embed files/coding_master/system.md
var coding_master string

//go:embed files/compare_and_contrast/system.md
var compare_and_contrast string

//go:embed files/create_5_sentence_summary/system.md
var create_5_sentence_summary string

//go:embed files/create_academic_paper/system.md
var create_academic_paper string

//go:embed files/create_ai_jobs_analysis/system.md
var create_ai_jobs_analysis string

//go:embed files/create_aphorisms/system.md
var create_aphorisms string

//go:embed files/create_art_prompt/system.md
var create_art_prompt string

//go:embed files/create_better_frame/system.md
var create_better_frame string

//go:embed files/create_coding_project/system.md
var create_coding_project string

//go:embed files/create_command/system.md
var create_command string

//go:embed files/create_cyber_summary/system.md
var create_cyber_summary string

//go:embed files/create_design_document/system.md
var create_design_document string

//go:embed files/create_diy/system.md
var create_diy string

//go:embed files/create_formal_email/system.md
var create_formal_email string

//go:embed files/create_git_diff_commit/system.md
var create_git_diff_commit string

//go:embed files/create_graph_from_input/system.md
var create_graph_from_input string

//go:embed files/create_hormozi_offer/system.md
var create_hormozi_offer string

//go:embed files/create_idea_compass/system.md
var create_idea_compass string

//go:embed files/create_investigation_visualization/system.md
var create_investigation_visualization string

//go:embed files/create_keynote/system.md
var create_keynote string

//go:embed files/create_logo/system.md
var create_logo string

//go:embed files/create_markmap_visualization/system.md
var create_markmap_visualization string

//go:embed files/create_mermaid_visualization/system.md
var create_mermaid_visualization string

//go:embed files/create_mermaid_visualization_for_github/system.md
var create_mermaid_visualization_for_github string

//go:embed files/create_micro_summary/system.md
var create_micro_summary string

//go:embed files/create_network_threat_landscape/system.md
var create_network_threat_landscape string

//go:embed files/create_newsletter_entry/system.md
var create_newsletter_entry string

//go:embed files/create_npc/system.md
var create_npc string

//go:embed files/create_pattern/system.md
var create_pattern string

//go:embed files/create_quiz/system.md
var create_quiz string

//go:embed files/create_reading_plan/system.md
var create_reading_plan string

//go:embed files/create_recursive_outline/system.md
var create_recursive_outline string

//go:embed files/create_report_finding/system.md
var create_report_finding string

//go:embed files/create_rpg_summary/system.md
var create_rpg_summary string

//go:embed files/create_security_update/system.md
var create_security_update string

//go:embed files/create_show_intro/system.md
var create_show_intro string

//go:embed files/create_sigma_rules/system.md
var create_sigma_rules string

//go:embed files/create_story_explanation/system.md
var create_story_explanation string

//go:embed files/create_stride_threat_model/system.md
var create_stride_threat_model string

//go:embed files/create_summary/system.md
var create_summary string

//go:embed files/create_tags/system.md
var create_tags string

//go:embed files/create_threat_scenarios/system.md
var create_threat_scenarios string

//go:embed files/create_ttrc_graph/system.md
var create_ttrc_graph string

//go:embed files/create_ttrc_narrative/system.md
var create_ttrc_narrative string

//go:embed files/create_upgrade_pack/system.md
var create_upgrade_pack string

//go:embed files/create_user_story/system.md
var create_user_story string

//go:embed files/create_video_chapters/system.md
var create_video_chapters string

//go:embed files/create_visualization/system.md
var create_visualization string

//go:embed files/dialog_with_socrates/system.md
var dialog_with_socrates string

//go:embed files/explain_code/system.md
var explain_code string

//go:embed files/explain_docs/system.md
var explain_docs string

//go:embed files/explain_math/system.md
var explain_math string

//go:embed files/explain_project/system.md
var explain_project string

//go:embed files/explain_terms/system.md
var explain_terms string

//go:embed files/export_data_as_csv/system.md
var export_data_as_csv string

//go:embed files/extract_algorithm_update_recommendations/system.md
var extract_algorithm_update_recommendations string

//go:embed files/extract_article_wisdom/system.md
var extract_article_wisdom string

//go:embed files/extract_book_ideas/system.md
var extract_book_ideas string

//go:embed files/extract_book_recommendations/system.md
var extract_book_recommendations string

//go:embed files/extract_business_ideas/system.md
var extract_business_ideas string

//go:embed files/extract_controversial_ideas/system.md
var extract_controversial_ideas string

//go:embed files/extract_core_message/system.md
var extract_core_message string

//go:embed files/extract_ctf_writeup/system.md
var extract_ctf_writeup string

//go:embed files/extract_extraordinary_claims/system.md
var extract_extraordinary_claims string

//go:embed files/extract_ideas/system.md
var extract_ideas string

//go:embed files/extract_insights/system.md
var extract_insights string

//go:embed files/extract_insights_dm/system.md
var extract_insights_dm string

//go:embed files/extract_instructions/system.md
var extract_instructions string

//go:embed files/extract_jokes/system.md
var extract_jokes string

//go:embed files/extract_latest_video/system.md
var extract_latest_video string

//go:embed files/extract_main_idea/system.md
var extract_main_idea string

//go:embed files/extract_most_redeeming_thing/system.md
var extract_most_redeeming_thing string

//go:embed files/extract_patterns/system.md
var extract_patterns string

//go:embed files/extract_poc/system.md
var extract_poc string

//go:embed files/extract_predictions/system.md
var extract_predictions string

//go:embed files/extract_primary_problem/system.md
var extract_primary_problem string

//go:embed files/extract_primary_solution/system.md
var extract_primary_solution string

//go:embed files/extract_product_features/system.md
var extract_product_features string

//go:embed files/extract_questions/system.md
var extract_questions string

//go:embed files/extract_recipe/system.md
var extract_recipe string

//go:embed files/extract_recommendations/system.md
var extract_recommendations string

//go:embed files/extract_references/system.md
var extract_references string

//go:embed files/extract_skills/system.md
var extract_skills string

//go:embed files/extract_song_meaning/system.md
var extract_song_meaning string

//go:embed files/extract_sponsors/system.md
var extract_sponsors string

//go:embed files/extract_videoid/system.md
var extract_videoid string

//go:embed files/extract_wisdom/system.md
var extract_wisdom string

//go:embed files/extract_wisdom_agents/system.md
var extract_wisdom_agents string

//go:embed files/extract_wisdom_dm/system.md
var extract_wisdom_dm string

//go:embed files/extract_wisdom_nometa/system.md
var extract_wisdom_nometa string

//go:embed files/find_hidden_message/system.md
var find_hidden_message string

//go:embed files/find_logical_fallacies/system.md
var find_logical_fallacies string

//go:embed files/get_wow_per_minute/system.md
var get_wow_per_minute string

//go:embed files/get_youtube_rss/system.md
var get_youtube_rss string

//go:embed files/identify_dsrp_distinctions/system.md
var identify_dsrp_distinctions string

//go:embed files/identify_dsrp_perspectives/system.md
var identify_dsrp_perspectives string

//go:embed files/identify_dsrp_relationships/system.md
var identify_dsrp_relationships string

//go:embed files/identify_dsrp_systems/system.md
var identify_dsrp_systems string

//go:embed files/identify_job_stories/system.md
var identify_job_stories string

//go:embed files/improve_academic_writing/system.md
var improve_academic_writing string

//go:embed files/improve_prompt/system.md
var improve_prompt string

//go:embed files/improve_report_finding/system.md
var improve_report_finding string

//go:embed files/improve_writing/system.md
var improve_writing string

//go:embed files/label_and_rate/system.md
var label_and_rate string

//go:embed files/md_callout/system.md
var md_callout string

//go:embed files/official_pattern_template/system.md
var official_pattern_template string

//go:embed files/prepare_7s_strategy/system.md
var prepare_7s_strategy string

//go:embed files/provide_guidance/system.md
var provide_guidance string

//go:embed files/rate_ai_response/system.md
var rate_ai_response string

//go:embed files/rate_ai_result/system.md
var rate_ai_result string

//go:embed files/rate_content/system.md
var rate_content string

//go:embed files/rate_value/system.md
var rate_value string

//go:embed files/raw_query/system.md
var raw_query string

//go:embed files/recommend_artists/system.md
var recommend_artists string

//go:embed files/recommend_pipeline_upgrades/system.md
var recommend_pipeline_upgrades string

//go:embed files/recommend_talkpanel_topics/system.md
var recommend_talkpanel_topics string

//go:embed files/refine_design_document/system.md
var refine_design_document string

//go:embed files/review_design/system.md
var review_design string

//go:embed files/show_fabric_options_markmap/system.md
var show_fabric_options_markmap string

//go:embed files/solve_with_cot/system.md
var solve_with_cot string

//go:embed files/suggest_pattern/system.md
var suggest_pattern string

//go:embed files/summarize/system.md
var summarize string

//go:embed files/summarize_debate/system.md
var summarize_debate string

//go:embed files/summarize_git_changes/system.md
var summarize_git_changes string

//go:embed files/summarize_git_diff/system.md
var summarize_git_diff string

//go:embed files/summarize_lecture/system.md
var summarize_lecture string

//go:embed files/summarize_legislation/system.md
var summarize_legislation string

//go:embed files/summarize_meeting/system.md
var summarize_meeting string

//go:embed files/summarize_micro/system.md
var summarize_micro string

//go:embed files/summarize_newsletter/system.md
var summarize_newsletter string

//go:embed files/summarize_paper/system.md
var summarize_paper string

//go:embed files/summarize_prompt/system.md
var summarize_prompt string

//go:embed files/summarize_pull-requests/system.md
var summarize_rpg_session string

//go:embed files/to_flashcards/system.md
var to_flashcards string

//go:embed files/transcribe_minutes/system.md
var transcribe_minutes string

//go:embed files/translate/system.md
var translate string

//go:embed files/tweet/system.md
var tweet string

//go:embed files/write_essay/system.md
var write_essay string

//go:embed files/write_hackerone_report/system.md
var write_hackerone_report string

//go:embed files/write_latex/system.md
var write_latex string

//go:embed files/write_micro_essay/system.md
var write_micro_essay string

//go:embed files/write_nuclei_template_rule/system.md
var write_nuclei_template_rule string

//go:embed files/write_semgrep_rule/system.md
var write_semgrep_rule string
